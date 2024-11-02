// Package server implements the Penne server
package server

import (
	"darvaza.org/sidecar/pkg/sidecar"
	"darvaza.org/sidecar/pkg/sidecar/horizon"
	"darvaza.org/slog"
	"darvaza.org/x/tls"

	"darvaza.org/penne/pkg/resolver"
)

// Server is a Penne server
type Server struct {
	cfg Config

	// sidecar
	sc *sidecar.Server
	// TLS
	tls tls.Store
	// horizons
	z horizon.Horizons
	// resolvers
	res map[string]*resolver.Resolver
	rd  map[string]slog.LogLevel
}

func (srv *Server) init() error {
	for _, fn := range []func() error{
		srv.initTLS,
		srv.initSidecar,
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
