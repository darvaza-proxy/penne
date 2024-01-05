package resolver

import (
	"crypto/tls"

	"github.com/miekg/dns"

	"darvaza.org/resolver/pkg/client"
	"darvaza.org/slog"
	"darvaza.org/slog/handlers/discard"
)

// Options contains information used to assemble all [Resolver]s.
type Options struct {
	Logger slog.Logger

	TLSConfig *tls.Config
}

// SetDefaults fills any gap in the [Options].
func (opts *Options) SetDefaults() {
	if opts.Logger == nil {
		opts.Logger = discard.New()
	}
}

// NewClient uses the [Options] to create a new [dns.Client].
func (opts *Options) NewClient(net string) client.Client {
	c := &dns.Client{
		Net:       net,
		TLSConfig: opts.TLSConfig,
	}

	switch net {
	case "tcp", "udp":
		c.TLSConfig = nil
	case "tcp+tls":
		if c.TLSConfig == nil {
			// not supported
			return nil
		}
	default:
		// not supported
		return nil
	}

	return c
}
