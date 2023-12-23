package server

import (
	"context"

	"darvaza.org/slog"
	"darvaza.org/slog/handlers/discard"

	"darvaza.org/penne/pkg/horizon"
	"darvaza.org/penne/pkg/resolver"
)

// Config describes how the Application will be assembled
type Config struct {
	Context context.Context `yaml:"-" toml:"-" json:"-"`
	Logger  slog.Logger     `yaml:"-" toml:"-" json:"-"`

	Horizons  []horizon.Config  `yaml:"horizons,omitempty" toml:",omitempty" json:",omitempty"`
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

	return nil
}
