package resolver

import (
	"darvaza.org/core"
	"darvaza.org/resolver"
	"darvaza.org/sidecar/pkg/sidecar/horizon"
	"darvaza.org/slog"
)

// Config describes a [Resolver].
type Config struct {
	Name string `yaml:"name"`
	Next string `yaml:"next,omitempty" toml:",omitempty" json:",omitempty"`

	// Debug indicates the requests passing through this [Resolver] should be logged or not.
	Debug bool `yaml:"debug,omitempty"        toml:",omitempty" json:",omitempty"`

	DisableAAAA bool     `yaml:"disable_aaaa,omitempty" toml:",omitempty" json:",omitempty"`
	Iterative   bool     `yaml:"iterative,omitempty"    toml:",omitempty" json:",omitempty"`
	Recursive   bool     `yaml:"recursive,omitempty"    toml:",omitempty" json:",omitempty"`
	Servers     []string `yaml:"servers,omitempty"      toml:",omitempty" json:",omitempty"`
	Suffixes    []string `yaml:"suffixes,omitempty"     toml:",omitempty" json:",omitempty"`

	Rewrites []RewriteConfig `yaml:"rewrite,omitempty" toml:",omitempty" json:",omitempty"`
}

// New creates a new [Resolver].
func (rc Config) New(next resolver.Exchanger, opts *Options) (*Resolver, error) {
	var err error

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

	switch {
	case rc.Iterative:
		err = rc.setupIterative(r, opts)
	case len(rc.Servers) > 0:
		err = rc.setupForwarder(r, opts)
	default:
		err = rc.setupChained(r, opts)
	}

	if err != nil {
		return nil, err
	}

	return r, nil
}

func (rc Config) setupForwarder(_ *Resolver, _ *Options) error {
	return &Error{
		Resolver: rc.Name,
		Reason:   "forwarder",
		Err:      core.ErrNotImplemented,
	}
}

func (rc Config) setupChained(_ *Resolver, _ *Options) error {
	return &Error{
		Resolver: rc.Name,
		Reason:   "chained resolver",
		Err:      core.ErrNotImplemented,
	}
}

// RewriteConfig describes an expression used to alter a request.
type RewriteConfig struct {
	From string `yaml:"from,omitempty" toml:",omitempty" json:",omitempty"`
	To   string `yaml:"to,omitempty" toml:",omitempty" json:",omitempty"`
}
