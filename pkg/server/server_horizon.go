package server

import (
	"darvaza.org/core"

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

func (*Server) initHorizons() error {
	return core.ErrNotImplemented
}
