package horizon

import (
	"context"

	"github.com/miekg/dns"

	"darvaza.org/resolver"
)

var _ resolver.Exchanger = (*Horizon)(nil)

// Exchange handles DNS requests passed from another [Horizon].
func (z *Horizon) Exchange(ctx context.Context, req *dns.Msg) (*dns.Msg, error) {
	var next resolver.Exchanger

	switch {
	case z.res != nil:
		// use explicit resolver
		next = z.res
	case z.next != nil:
		// hand-over to the next Horizon
		next = z.next
	default:
		// EOL
		next = z.nextE
	}

	return next.Exchange(ctx, req)
}

// HorizonExchange handles DNS requests directly from the [dns.Server] when the
// client belongs in the range.
//
// A Horizon that acts as entry point has to make sure security constraints
// are checked.
func (z *Horizon) HorizonExchange(ctx context.Context, req *dns.Msg) (*dns.Msg, error) {
	// TODO: replace EDNS0 SUBNET when forwarding isn't allowed
	return z.Exchange(ctx, req)
}
