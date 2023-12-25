package resolver

import (
	"darvaza.org/slog"
	"darvaza.org/slog/handlers/discard"
)

// Options contains information used to assemble all [Resolver]s.
type Options struct {
	Logger slog.Logger
}

// SetDefaults fills any gap in the [Options].
func (opts *Options) SetDefaults() {
	if opts.Logger == nil {
		opts.Logger = discard.New()
	}
}
