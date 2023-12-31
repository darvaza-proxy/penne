package server

import (
	"context"
	"time"

	"github.com/amery/defaults"

	"darvaza.org/slog"
	"darvaza.org/slog/handlers/discard"

	"darvaza.org/penne/pkg/horizon"
	"darvaza.org/penne/pkg/resolver"
)

// Config describes how the Application will be assembled
type Config struct {
	Context context.Context `yaml:"-" toml:"-" json:"-"`
	Logger  slog.Logger     `yaml:"-" toml:"-" json:"-"`

	Name    string `yaml:"name"    default:"localhost"`
	Version string `yaml:"version" default:"unspecified"`
	Authors string `yaml:"authors" default:"JPI Technologies <oss@jpi.io>"`

	// DisableCHAOS makes the DNS server respond with an empty success instead of giving
	// away software information.
	DisableCHAOS bool `yaml:"disable_chaos,omitempty" toml:",omitempty" json:",omitempty"`

	// ExchangeTimeout indicates the deadline to be used on DNS requests
	ExchangeTimeout time.Duration `yaml:"exchange_timeout" default:"5s"`

	Horizons  []horizon.Config  `yaml:"horizons,omitempty"  toml:",omitempty" json:",omitempty"`
	Resolvers []resolver.Config `yaml:"resolvers,omitempty" toml:",omitempty" json:",omitempty"`
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
	return defaults.Set(cfg)
}
