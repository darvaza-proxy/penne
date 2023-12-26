package resolver

import (
	"darvaza.org/resolver"
	"darvaza.org/sidecar/pkg/sidecar/horizon"
	"darvaza.org/slog"
)

// Config describes a [Resolver].
type Config struct {
	Name string `yaml:"name"`
	Next string `yaml:"next,omitempty" toml:",omitempty" json:",omitempty"`

	DisableAAAA bool     `yaml:"disable_aaaa,omitempty" toml:",omitempty" json:",omitempty"`
	Iterative   bool     `yaml:"iterative,omitempty"    toml:",omitempty" json:",omitempty"`
	Recursive   bool     `yaml:"recursive,omitempty"    toml:",omitempty" json:",omitempty"`
	Servers     []string `yaml:"servers,omitempty"      toml:",omitempty" json:",omitempty"`
	Suffixes    []string `yaml:"suffixes,omitempty"     toml:",omitempty" json:",omitempty"`

	Rewrites []RewriteConfig `yaml:"rewrite,omitempty" toml:",omitempty" json:",omitempty"`
}

// New creates a new [Resolver].
func (rc Config) New(next resolver.Exchanger, opts *Options) (*Resolver, error) {
	if opts == nil {
		opts = new(Options)
	}
	opts.SetDefaults()

	r := &Resolver{
		debug:    make(map[string]slog.LogLevel),
		log:      opts.Logger,
		name:     rc.Name,
		suffixes: rc.Suffixes,

		Next:      next,
		Exchanger: resolver.ExchangerFunc(horizon.ForbiddenExchange),
	}

	// TODO: set them up to do something
	return r, nil
}

// RewriteConfig describes an expression used to alter a request.
type RewriteConfig struct {
	From string `yaml:"from,omitempty" toml:",omitempty" json:",omitempty"`
	To   string `yaml:"to,omitempty" toml:",omitempty" json:",omitempty"`
}
