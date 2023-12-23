package server

import (
	"darvaza.org/penne/pkg/resolver"
)

func defaultResolvers() []resolver.Config {
	return []resolver.Config{
		{
			Name:      "root",
			Iterative: true,
		},
	}
}

func (srv *Server) initResolvers() error {
	_, _, err := resolver.MakeResolvers(srv.cfg.Resolvers, srv.cfg.Logger)
	return err
}
