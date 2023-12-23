package server

import (
	"darvaza.org/penne/pkg/horizon"
)

func defaultHorizons() []horizon.Config {
	return []horizon.Config{
		{
			Name:     "any",
			Resolver: "root",
		},
	}
}

func (srv *Server) initHorizons() error {
	_, _, err := horizon.MakeHorizons(srv.cfg.Horizons, nil)
	return err
}
