package server

import (
	"darvaza.org/core"

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

func (*Server) initResolvers() error {
	return core.ErrNotImplemented
}
