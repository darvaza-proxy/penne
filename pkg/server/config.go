package server

import (
	"context"

	"darvaza.org/slog"
	"darvaza.org/slog/handlers/discard"
)

// Config describes how the Application will be assembled
type Config struct {
	Context context.Context `yaml:"-" toml:"-" json:"-"`
	Logger  slog.Logger     `yaml:"-" toml:"-" json:"-"`
}

// SetDefaults fills gaps in the Config
func (cfg *Config) SetDefaults() error {
	if cfg.Context == nil {
		cfg.Context = context.Background()
	}

	if cfg.Logger == nil {
		cfg.Logger = discard.New()
	}

	return nil
}
