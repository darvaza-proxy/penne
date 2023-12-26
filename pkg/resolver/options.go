package resolver

import (
	"github.com/miekg/dns"

	"darvaza.org/resolver/pkg/client"
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

// NewClient uses the [Options] to create a new [dns.Client].
func (*Options) NewClient(net string) client.Client {
	return &dns.Client{Net: net}
}
