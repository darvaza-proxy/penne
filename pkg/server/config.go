package server

import (
	"context"
	"net/netip"
	"time"

	"darvaza.org/sidecar/pkg/sidecar"
	"darvaza.org/sidecar/pkg/sidecar/store"
	"darvaza.org/slog"
	"darvaza.org/slog/handlers/discard"
	"darvaza.org/x/config"

	"darvaza.org/penne/pkg/horizon"
	"darvaza.org/penne/pkg/resolver"
)

// Config describes how the Application will be assembled
type Config struct {
	Context context.Context `yaml:"-" toml:"-" json:"-"`
	Logger  slog.Logger     `yaml:"-" toml:"-" json:"-"`

	Name    string `default:"localhost"`
	Version string `default:"unspecified"`
	Authors string `default:"JPI Technologies <oss@jpi.io>"`

	// DisableCHAOS makes the DNS server respond with an empty success instead of giving
	// away software information.
	DisableCHAOS bool `yaml:"disable_chaos,omitempty" toml:",omitempty" json:",omitempty"`

	// ExchangeTimeout indicates the deadline to be used on DNS requests
	ExchangeTimeout time.Duration `yaml:"exchange_timeout" default:"5s"`

	// Supervision provides details for graceful upgrades and restarts
	Supervision sidecar.SupervisionConfig `yaml:"supervision,omitempty" toml:",omitempty" json:",omitempty"`

	// Listen provides details on the addresses and ports to handle
	Listen ListenConfig `yaml:",omitempty" toml:",omitempty" json:",omitempty"`

	// TLS contains instructions to assemble the TLS store.
	// TODO: allow ACME
	TLS store.Config `yaml:",omitempty" toml:",omitempty" json:",omitempty"`

	Horizons  []horizon.Config  `yaml:",omitempty" toml:",omitempty" json:",omitempty"`
	Resolvers []resolver.Config `yaml:",omitempty" toml:",omitempty" json:",omitempty"`
}

// ListenConfig describes the addresses and ports this [Server] will listen
type ListenConfig struct {
	Interfaces []string     `yaml:",omitempty" toml:",omitempty" json:",omitempty"`
	Addresses  []netip.Addr `yaml:",omitempty"  toml:",omitempty" json:",omitempty"`

	DNS   uint16 `yaml:"dns_port"   default:"53"`
	DoT   uint16 `yaml:"dot_port"   default:"853"`
	HTTP  uint16 `yaml:"http_port"  default:"80"`
	HTTPS uint16 `yaml:"https_port" default:"443"`

	// DisableHTTP disables support for plain/unencrypted HTTP of the monitoring interface
	DisableHTTP bool `yaml:"disable_http,omitempty" toml:",omitempty" json:",omitempty"`
}

// SetDefaults fills gaps in the Config
func (cfg *Config) SetDefaults() error {
	if cfg.Context == nil {
		cfg.Context = context.Background()
	}

	if cfg.Logger == nil {
		cfg.Logger = discard.New()
	}

	if len(cfg.Horizons) == 0 {
		cfg.Horizons = defaultHorizons()
	}

	if len(cfg.Resolvers) == 0 {
		cfg.Resolvers = defaultResolvers()
	}

	// and the rest
	return config.Set(cfg)
}
