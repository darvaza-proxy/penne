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

		Addresses: sidecar.BindConfig{
			Interfaces: cfg.Listen.Interfaces,
			Addresses:  cfg.Listen.Addresses,
		},

		HTTP: sidecar.HTTPConfig{
			Port:           cfg.Listen.HTTPS,
			PortInsecure:   cfg.Listen.HTTP,
			EnableInsecure: !cfg.Listen.DisableHTTP,
		},

		DNS: sidecar.DNSConfig{
			Enabled: true,
			Port:    cfg.Listen.DNS,
			TLSPort: cfg.Listen.DoT,
		},
	}

	if err := scc.SetDefaults(); err != nil {
		return nil, err
	}

	return scc, nil
}
