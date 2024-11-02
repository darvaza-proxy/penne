package server

import (
	"darvaza.org/sidecar/pkg/sidecar"
	"darvaza.org/x/tls"
)

func (cfg *Config) export(s tls.Store) (*sidecar.Config, error) {
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
