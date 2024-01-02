package server

import (
	"darvaza.org/darvaza/shared/storage"
	"darvaza.org/sidecar/pkg/sidecar"
)

func (cfg *Config) export(s storage.Store) (*sidecar.Config, error) {
	scc := &sidecar.Config{
		Context: cfg.Context,
		Logger:  cfg.Logger,
		Store:   s,

		Name: cfg.Name,

		Supervision: cfg.Supervision,
	}

	if err := scc.SetDefaults(); err != nil {
		return nil, err
	}

	return scc, nil
}
