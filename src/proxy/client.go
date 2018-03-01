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
    "io"
	"net"
    "os"
    "time"
    "github.com/darrylwest/go-unique/unique"
)

// Client the client object, created for each new request
type Client struct {
    cfg *Config
    id  string
    created time.Time
    request *bytes.Buffer
    response *bytes.Buffer
}

// NewClient create a new client with config, id and date created
func NewClient(cfg *Config) *Client {
    return &Client{
        cfg: cfg,
        id: unique.CreateULID(),
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

// ReadRequest reads the entire request and stores in client.request
func (client Client) ReadRequest(dst io.Writer, src io.Reader) error {
    size := 32 * 1024
    buf := make([]byte, size)
    copied := 0
    var err error
    loopCount := 0

    for {
        nr, er := src.Read(buf)
        loopCount++

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

            if loopCount > 1 {
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

    return err
}

// SendResponse sends the response back to the original requestor
func (client Client) SendResponse(dst io.Writer) error {
    log.Info("send response...")
    payload := []byte(`{"status":"ok"}` + "\r\n")

    n, err := dst.Write(payload)
    log.Info("%d bytes written...", n)

    return err
}

func (client Client) handleRequest(sock net.Conn) error {
    defer sock.Close()
    log.Info("handle request: %s %s", client.id, client.created.Format(time.RFC3339))

    // read the request in full
    sock.SetReadDeadline(time.Now().Add(20 * time.Second))

    readComplete := make(chan bool)
    go func() {
        // filename := fmt.Sprintf("data/%s-request.log", client.id)
        client.request = new(bytes.Buffer)
        err := client.ReadRequest(client.request, sock)
        if err != nil {
            log.Warn("%s", err)
        }

        log.Info("%s", client.request.String())
        readComplete <- true
    }()

    <- readComplete

    err := client.SendResponse(sock)
    if err != nil {
        log.Error("write : %s", err)
    }

    return err
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

