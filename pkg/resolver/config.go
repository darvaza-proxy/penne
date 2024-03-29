package resolver

import (
	"fmt"

	"darvaza.org/core"
	"darvaza.org/resolver"
	"darvaza.org/slog"

	"darvaza.org/penne/pkg/suffix"
)

// Config describes a [Resolver].
type Config struct {
	// Name is the unique name of this [Resolver]
	Name string `yaml:""`
	// Next is the name of the resolver to use if the Suffixes restriction
	// isn't satisfied.
	Next string `yaml:",omitempty" toml:",omitempty" json:",omitempty"`

	// Debug indicates the requests passing through this [Resolver] should be logged or not.
	Debug bool `yaml:",omitempty" toml:",omitempty" json:",omitempty"`

	// OmitSubNet indicates requests reaching out to remote servers should omit
	// EDNS0 SUBNET information.
	OmitSubNet bool `yaml:"omit_subnet,omitempty" toml:",omitempty" json:",omitempty"`

	// DisableAAAA indicates that this [Resolver] will discard AAAA entries
	DisableAAAA bool `yaml:"disable_aaaa,omitempty" toml:",omitempty" json:",omitempty"`

	// Iterative indicates that this [Resolver] will go straight to the DNS
	// root servers and ask the authoritative servers for the answers.
	Iterative bool `yaml:",omitempty" toml:",omitempty" json:",omitempty"`
	// IterativeMaxRR indicates the maximum number of glue records the Iterative Resolver
	// will cache.
	IterativeMaxRR uint `yaml:",omitempty" toml:",omitempty" json:",omitempty"`

	// Recursive indicates that this [Resolver] will ask servers to perform
	// recursive lookups.
	Recursive bool `yaml:",omitempty" toml:",omitempty" json:",omitempty"`

	// Servers is a list of DNS servers to use for forwarding or iterative resolution.
	// If this [Resolver] is designated as iterative and no servers are provided,
	// a built-in list of root DNS servers will be used.
	Servers []string `yaml:",omitempty" toml:",omitempty" json:",omitempty"`

	// Workers indicates how many parallel outward connections we allow
	// from this resolver. It applies to iterative and forwarding resolvers.
	Workers uint `yaml:",omitempty" toml:",omitempty" json:",omitempty"`

	// Suffixes indicate what domains will this [Resolver] handle. Globbing patterns allowed.
	Suffixes []string `yaml:",omitempty" toml:",omitempty" json:",omitempty"`

	// Rewrites is a list of query name rewrites to be done by this [Resolver].
	Rewrites []RewriteConfig `yaml:",omitempty" toml:",omitempty" json:",omitempty"`
}

// SetDefaults fills gaps in the [Config].
func (rc *Config) SetDefaults() {
	if rc.Iterative {
		rc.Recursive = false

		if rc.IterativeMaxRR == 0 {
			rc.IterativeMaxRR = DefaultIteratorMaxRR
		}
	} else {
		rc.IterativeMaxRR = 0
	}
}

// New creates a new [Resolver].
func (rc Config) New(next resolver.Exchanger, opts *Options) (*Resolver, error) {
	var err error

	if opts == nil {
		opts = new(Options)
	}
	opts.SetDefaults()

	suffixes, err := suffix.MakeSuffixes(rc.Suffixes)
	if err != nil {
		return nil, err
	}

	r := &Resolver{
		debug:    make(map[string]slog.LogLevel),
		log:      opts.Logger,
		name:     rc.Name,
		suffixes: suffixes,
		next:     next,
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

func (rc Config) setupChained(_ *Resolver, _ *Options) error {
	return rc.WrapError(core.ErrNotImplemented, "chained resolver")
}

// NewError creates a new [Error] using the [Config]'s name.
func (rc Config) NewError(format string, args ...any) *Error {
	var s string
	if len(args) > 0 {
		s = fmt.Sprintf(format, args...)
	} else {
		s = format
	}

	return &Error{
		Resolver: rc.Name,
		Reason:   s,
	}
}

// WrapError creates a new [Error] using the [Config]'s name, wrapping
// another error.
func (rc Config) WrapError(err error, format string, args ...any) error {
	e := rc.NewError(format, args...)
	e.Err = err
	return e
}

// RewriteConfig describes an expression used to alter a request.
type RewriteConfig struct {
	// From is a globbing pattern to match and capture.
	From string `yaml:",omitempty" toml:",omitempty" json:",omitempty"`
	// To is the rewrite template for entries that match the `From` pattern.
	To string `yaml:",omitempty" toml:",omitempty" json:",omitempty"`
	// Final indicates that entries matching this From shouldn't continue
	// to the next rewrite rule.
	Final bool
}
