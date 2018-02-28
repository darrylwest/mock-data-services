//
// service - define the routes and start the service
//
// @author darryl.west@ebay.com
// @created 2018-02-27 16:22:25
//

package proxy

import (
	"fmt"
	"net"
)

// Service - the service struct
type Service struct {
	cfg *Config
}

// NewService create a new service by passing in config
func NewService(cfg *Config) (*Service, error) {
	svc := Service{
		cfg: cfg,
	}

	return &svc, nil
}

// Start start the admin listener and event loop
func (svc Service) Start() error {
	log.Info("start the hub service...")
	// cfg := svc.cfg

	// create and open the registry database

	// start the listener
	if err := svc.startServer(); err != nil {
		return err
	}

	return nil
}

func (svc Service) startServer() error {
	cfg := svc.cfg
    port := cfg.Port

	host := fmt.Sprintf(":%d", cfg.Port)
	log.Info("start listening on port %s", host)

	// open the listening socket
    listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
		log.Error("could not start server on %d: %v", port, err)
		return err
    }

	log.Info("proxy listening on %d, proxy to %s\n", port, cfg.Target)

    for {
        sock, err := listener.Accept()
		if err != nil {
			log.Error("could not accept client connection", err)
			return err
		}

        go func() {
            defer sock.Close()
            client := NewClient(cfg)

            if err = client.handleRequest(sock); err != nil {
                log.Error("client %s error: %s", client.id, err);
            }

            log.Info("%s closed...", client.id)
        }()
    }
}
