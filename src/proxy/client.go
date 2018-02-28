//
// client - define the routes and start the service
//
// @author darryl.west@ebay.com
// @created 2018-02-28 16:22:25
//

package proxy

import (
	"net"
    "time"
    "github.com/darrylwest/go-unique/unique"
)

// Client the client object, created for each new request
type Client struct {
    cfg *Config
    id  string
    created time.Time
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

func (client Client) handleRequest(sock net.Conn) error {
    return nil
}

