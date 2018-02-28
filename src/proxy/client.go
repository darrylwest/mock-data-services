//
// client - define the routes and start the service
//
// @author darryl.west@ebay.com
// @created 2018-02-28 16:22:25
//

package proxy

import (
    "io"
	"net"
    "time"
    "github.com/darrylwest/go-unique/unique"
)

// Client the client object, created for each new request
type Client struct {
    cfg *Config
    id  string
    created time.Time
    request []byte
    response []byte
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
func (client *Client) ReadRequest(src io.Reader) error {
    size := 32 * 1024
    buf := make([]byte, size)
    var err error

    for {
        n, e := src.Read(buf)
        if n > 0 {
            client.request = append(client.request, buf[0:n]...)
        }

        if e != nil {
            if e != io.EOF {
                err = e
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
    sock.SetReadDeadline(time.Now().Add(1 * time.Millisecond))
    err := client.ReadRequest(sock)

    log.Info("request: %s", client.request)

    err = client.SendResponse(sock)

    return err
}

