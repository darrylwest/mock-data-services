//
// client - define the routes and start the service
//
// @author darryl.west@ebay.com
// @created 2018-02-28 16:22:25
//

package proxy

import (
	"bytes"
	"fmt"
	"github.com/darrylwest/go-unique/unique"
	"io"
	"net"
	"os"
    "strings"
    "strconv"
	"time"
)

const mockresp = `HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 01 Mar 2018 15:03:07 GMT
Content-Length: 262

{"status":"ok","ts":1519916587296,"version":"1.0","webStatus":{"agent":"curl/7.54.0","host":"127.0.0.1:3400","path":"/status","pid":9819,"proto":"HTTP/1.1","remoteAddr":"127.0.0.1:38396","version":"0.90.128","xForwardedFor":"73.158.29.165","xForwardedProto":""}}
`

// ClientRequest the parsed client request 
type ClientRequest struct {
    method string
    uri    string
    size   int
    duration time.Duration
}

// ClientResponse holds status and content length
type ClientResponse struct {
    status string
    size   int
}

// Client the client object, created for each new request
type Client struct {
	cfg      *Config
	id       string
	created  time.Time
	request  *bytes.Buffer
	response *bytes.Buffer
}

// NewClient create a new client with config, id and date created
func NewClient(cfg *Config) *Client {
	return &Client{
		cfg:     cfg,
		id:      unique.CreateULID(),
		created: time.Now(),
	}
}

// GetID returns the client's id
func (client Client) GetID() string {
	return client.id
}

// GetCreatedAt returns the date/time that this client was created
func (client Client) GetCreatedAt() time.Time {
	return client.created
}

// ParseContentLength parse the content length header
func (client Client) ParseContentLength(line string) (int, error) {
    cols := strings.Split(line, ":")
    sz := strings.Trim(cols[1], " \r")

    val, err := strconv.Atoi(string(sz))
    if err != nil {
        log.Error("parse content length error: %s : %s", sz, err)
    }

    return val, err
}

// ParseResponse return the status and content length
func (client Client) ParseResponse(buf []byte) *ClientResponse {
    lines := bytes.Split(buf, []byte("\n"))
    resp := ClientResponse{}

    for idx, line := range lines {
        if idx == 0 {
            resp.status = string(line)

            continue
        }

        // parse the content length
        if len(line) < 80 && bytes.HasPrefix(bytes.ToUpper(line), []byte("CONTENT-LENGTH:")) {
            log.Info("parse header: %s", line)
            if val, err := client.ParseContentLength(string(line)); err == nil {
                resp.size = val
            }

            break
        }

    }

    return &resp
}

// ParseRequest parses the method, uri and content length
func (client Client) ParseRequest(buf []byte) *ClientRequest {
    lines := bytes.Split(buf, []byte("\n"))

    req := ClientRequest{}

    for idx, line := range lines {
        if idx == 0 {
            cols := bytes.Split(line, []byte(" "))
            req.method = string(cols[0])
            req.uri = string(cols[1])

            if req.method == "GET" {
                break
            }

            continue
        }

        // parse the content length
        if len(line) < 80 && bytes.HasPrefix(bytes.ToUpper(line), []byte("CONTENT-LENGTH:")) {
            log.Info("parse header: %s", line)
            if val, err := client.ParseContentLength(string(line)); err == nil {
                req.size = val
            }

            break
        }
    }

    log.Info("line count: %d, method: %s, uri: %s, size: %d", len(lines), req.method, req.uri, req.size)

    return &req
}

// ReadRequest reads the entire request and stores in client.request
func (client Client) ReadRequest(dst io.Writer, src io.Reader) (*ClientRequest, error) {
	size := 16 * 1024
	buf := make([]byte, size)
	copied := 0
	var err error
    var req *ClientRequest

	for {
		nr, er := src.Read(buf)

		if nr > 0 {
			nw, er := dst.Write(buf[0:nr])
			if nw > 0 {
				copied += nw
			}
			if er != nil {
				err = er
				break
			}
			if nw != nr {
				err = fmt.Errorf("number of bytes written does not match")
				break
			}

			if req == nil {
				log.Info("parse %d bytes", nr)
                req = client.ParseRequest(buf[0:nr])
			}

            if req.method == "GET" || copied >= req.size {
                break
            }
		}

		if er != nil {
			if er != io.EOF {
				err = er
			}

			break
		}
	}

	return req, err
}

// String show the volues for stats file
func (cr ClientRequest) String() string {
    return fmt.Sprintf("method: %s\nuri: %s\nbytes: %d\nduration: %s\n", cr.method, cr.uri, cr.size, cr.duration)
}

// ProxyRequest send the request bytes to the target
func (client *Client) ProxyRequest(dst, sock net.Conn) error {
    n, err := sock.Write(client.request.Bytes())
    log.Info("%d bytes of %d written to target", n, client.request.Len())

    if err != nil && err != io.EOF {
        return err
    }

    var resp *ClientResponse
    client.response = new(bytes.Buffer)
    buf := make([]byte, client.cfg.BufSize * 1024)
    written := 0
    for {
        nr, er := sock.Read(buf)
        if nr > 0 {
            log.Info("read from target %d bytes", nr)
            nw, ew := dst.Write(buf[0:nr])
            if nw > 0 {
                client.response.Write(buf[0:nr])
                written += nw
            }

            if ew != nil {
                err = ew
                break
            }

            if nr != nw {
                err = fmt.Errorf("short write to client: %d != %d", nr, nw)
                break
            }

            if resp == nil {
                resp = client.ParseResponse(buf[0:nr])
            }

            if written >= resp.size {
                break
            }
        }

        if er != nil {
            if er != io.EOF {
                err = er
                break
            }
        }

    }

    return err
}

// SendResponse sends the response back to the original requestor
func (client Client) SendResponse(dst io.Writer, resp []byte) error {
	log.Info("send response, size: %d", len(resp))

	n, err := dst.Write(resp)
	log.Info("%d bytes written...", n)

	return err
}

func (client Client) openTarget() (net.Conn, error) {
    cfg := client.cfg
    var conn net.Conn
    if cfg.Bypass {
        return conn, nil
    }

    log.Info("open the target: %s", cfg.Target)
    conn, err := net.Dial("tcp", cfg.Target)

    if err == nil {
        conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
        conn.SetReadDeadline(time.Now().Add(10 * time.Second))
    }

    return conn, err
}

func (client *Client) handleRequest(sock net.Conn) error {
	defer sock.Close()
	log.Info("handle request: %s %s", client.id, client.created.Format(time.RFC3339))

	// read the request in full
	sock.SetReadDeadline(time.Now().Add(20 * time.Second))
	sock.SetWriteDeadline(time.Now().Add(10 * time.Second))

	readComplete := make(chan *ClientRequest)
	go func() {
		client.request = new(bytes.Buffer)
		req, err := client.ReadRequest(client.request, sock)
		if err != nil {
			log.Warn("%s", err)
		}

		readComplete <- req

        req.duration = time.Now().Sub(client.created)

		log.Info("client request size: %d, content-length: %d, read time: %s", client.request.Len(), req.size, req.duration)

        client.writeFile(client.GetRequestFilename(), client.request.Bytes())
        client.writeFile(client.GetStatsFilename(), []byte(req.String()))
	}()

    targetOpen := false
    target, err := client.openTarget()
    if err != nil {
        log.Error("error opening target: %s", err)
        defer target.Close()
    } else {
        targetOpen = true
    }

	clientRequest := <-readComplete

    if targetOpen {
        log.Info("write response to the target...")
        err = client.ProxyRequest(sock, target)
        if err != nil {
            log.Error("error proxying request: %s", err)
        }
        log.Info("response size: %d", client.response.Len())
    } else {
        client.response = client.GetMockResponse(clientRequest)
        err = client.SendResponse(sock, client.response.Bytes())
        if err != nil {
            log.Error("write : %s", err)
        }
    }

    // now write to the response log
    if er := client.writeFile(client.GetResponseFilename(), client.response.Bytes()); er != nil {
        log.Error("writing response file: %s", er)
    }

	return err
}

// GetMockResponse reads the registry based on method/uri and returns a mock
func (client Client) GetMockResponse(req *ClientRequest) *bytes.Buffer {
    return bytes.NewBuffer([]byte(mockresp))
}

func (client Client) writeFile(filename string, buf []byte) error {
	log.Info("write request to %s", filename)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.Write(buf); err != nil {
		return err
	}

	_, err = f.WriteString("\n")
	return err
}

// GetRequestFilename returns the request filename for this client
func (client Client) GetRequestFilename() string {
    return fmt.Sprintf("data/%s-request.log", client.id)
}

// GetResponseFilename returns the stats filename for this client
func (client Client) GetResponseFilename() string {
    return fmt.Sprintf("data/%s-response.log", client.id)
}

// GetStatsFilename returns the stats filename for this client
func (client Client) GetStatsFilename() string {
    return fmt.Sprintf("data/%s-stats.log", client.id)
}
