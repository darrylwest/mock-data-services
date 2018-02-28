//
// service - define the routes and start the service
//
// @author darryl.west@ebay.com
// @created 2018-02-27 16:22:25
//

package proxy

import (
	"fmt"
	// "net/http"
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

	host := fmt.Sprintf(":%d", cfg.Port)
	log.Info("start listening on port %s", host)

	// open the listening socket

	return nil
}
