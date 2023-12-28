package server

import (
	"darvaza.org/darvaza/shared/storage"
	"darvaza.org/sidecar/pkg/sidecar"
)

func (cfg *Config) export(s storage.Store) (*sidecar.Config, error) {
	addrs := make([]string, 0, len(cfg.Listen.Addresses))
	for _, addr := range cfg.Listen.Addresses {
		addrs = append(addrs, addr.String())
	}

	scc := &sidecar.Config{
		Context: cfg.Context,
		Logger:  cfg.Logger,
		Store:   s,

		Name: cfg.Name,

		Supervision: cfg.Supervision,

		Addresses: sidecar.BindConfig{
			Interfaces: cfg.Listen.Interfaces,
			Addresses:  addrs,
		},

		HTTP: sidecar.HTTPConfig{
			Port:           cfg.Listen.HTTPS,
			PortInsecure:   cfg.Listen.HTTP,
			EnableInsecure: !cfg.Listen.DisableHTTP,
		},
	}

	if err := scc.SetDefaults(); err != nil {
		return nil, err
	}

	return scc, nil
}
