package server

import (
	"darvaza.org/resolver/pkg/reflect"

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
	// prepare srv.z
	srv.z.ContextKey = horizon.NewContextKey("penne.horizon.match")
	srv.z.ExchangeTimeout = srv.cfg.ExchangeTimeout
	srv.z.ExchangeContext = reflect.WithEnabledFunc(srv.cfg.Context, srv.reflectEnabled)

	_, _, err := horizon.MakeHorizons(srv.cfg.Horizons, nil)
	return err
}
