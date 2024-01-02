// Package server implements the Penne server
package server

import (
	"darvaza.org/darvaza/shared/storage"
	"darvaza.org/sidecar/pkg/sidecar/horizon"
)

// Server is a Penne server
type Server struct {
	cfg Config

	// TLS
	tls storage.Store
	// horizons
	z horizon.Horizons
}

func (srv *Server) init() error {
	for _, fn := range []func() error{
		srv.initTLS,
		srv.initResolvers,
		srv.initHorizons,
	} {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

// New creates a new [Server] based on the given [Config]
func New(cfg *Config) (*Server, error) {
	if cfg == nil {
		cfg = new(Config)
	}

	if err := cfg.SetDefaults(); err != nil {
		return nil, err
	}

	srv := &Server{
		cfg: *cfg,
	}

	if err := srv.init(); err != nil {
		return nil, err
	}

	return srv, nil
}
