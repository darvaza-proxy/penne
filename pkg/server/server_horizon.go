package server

import (
	"net/http"

	"darvaza.org/resolver"
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

	// build horizons
	names, m, err := horizon.MakeHorizons(srv.cfg.Horizons, nil)
	if err != nil {
		return err
	}

	return srv.assembleHorizons(names, m)
}

func (srv *Server) assembleHorizons(names []string, m map[string]*horizon.Horizon) error {
	var h http.Handler
	var e resolver.Exchanger

	// TODO: set h to the handler of our web interface

	// preserve original order
	for _, name := range names {
		z := m[name]

		// create horizon.Horizon
		zp := z.New(h, e)
		if err := srv.z.Append(zp); err != nil {
			return err
		}
	}

	return nil
}
