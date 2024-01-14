package horizon

import (
	"context"

	"github.com/miekg/dns"

	"darvaza.org/core"
	"darvaza.org/resolver"
	"darvaza.org/resolver/pkg/exdns"
	"darvaza.org/sidecar/pkg/sidecar/horizon"
)

var _ resolver.Exchanger = (*Horizon)(nil)

// Exchange handles DNS requests passed from another [Horizon].
func (z *Horizon) Exchange(ctx context.Context, req *dns.Msg) (*dns.Msg, error) {
	if e := z.res; e != nil {
		// use explicit resolver with trampoline
		ctx = dnsNextKey.WithValue(ctx, z.nextExchanger)
		return e.Exchange(ctx, req)
	}

	return z.nextExchanger(ctx, req)
}

func (z *Horizon) nextExchanger(ctx context.Context, req *dns.Msg) (*dns.Msg, error) {
	var next resolver.Exchanger

	if z.next != nil {
		// next Horizon
		next = z.next
	} else {
		// EOL
		next = z.nextE
	}

	return next.Exchange(ctx, req)
}

// dnsNextEx is injected as the fallback on the Resolvers
// so they can continue to the next Horizon if they don't have
// anything else defined already.
func dnsNextEx(ctx context.Context, req *dns.Msg) (*dns.Msg, error) {
	e, ok := dnsNextKey.Get(ctx)
	if ok {
		return e.Exchange(ctx, req)
	}

	return horizon.ForbiddenExchange(ctx, req)
}

var (
	dnsNextKey       = core.NewContextKey[resolver.ExchangerFunc]("penne.horizon.resolver")
	dnsNextExchanger = resolver.ExchangerFunc(dnsNextEx)
)

// HorizonExchange handles DNS requests directly from the [dns.Server] when the
// client belongs in the range.
//
// A Horizon that acts as entry point has to make sure security constraints
// are checked.
func (z *Horizon) HorizonExchange(ctx context.Context, req *dns.Msg) (*dns.Msg, error) {
	var original = req

	switch {
	case !z.allowForwarding:
		if req2, ok := z.replaceEDNS0SUBNET(ctx, req); ok {
			req = req2
		}
	case !containsEDNS0SUBNET(req):
		// add EDNS0 SUBNET data based on the horizon.Match
		// if none is included already
		if o, ok := z.newEDNS0SUBNET(ctx); ok {
			// add EDNS0 SUBNET data based on the horizon.Match
			req2 := req.Copy()
			req2.Id = dns.Id()
			req2.Extra = append(req2.Extra, o)
			req = req2
		}
	}

	resp, err := z.Exchange(ctx, req)
	return exdns.RestoreReturn(original, resp, err)
}

func (z *Horizon) replaceEDNS0SUBNET(ctx context.Context, req *dns.Msg) (*dns.Msg, bool) {
	req2 := req.Copy()

	// remove EDNS0 SUBNET data
	altered := removeEDNS0SUBNET(req2)

	if o, ok := z.newEDNS0SUBNET(ctx); ok {
		// add EDNS0 SUBNET data based on the horizon.Match
		req2.Extra = append(req2.Extra, o)
		altered = true
	}

	if altered {
		// new request
		req2.Id = dns.Id()
		return req2, true
	}

	return nil, false
}

func (z *Horizon) newEDNS0SUBNET(ctx context.Context) (dns.RR, bool) {
	m, ok := z.ctxKey.Get(ctx)
	if !ok {
		// no horizon.Match data
		return nil, false
	}

	bits := m.CIDR.Bits()
	if bits == 0 {
		// don't add entry for /0
		return nil, false
	}

	addr := m.CIDR.Addr()
	family := core.IIf(addr.Is6(), 2, 1)

	// EDNS0 SUBNET
	e := new(dns.EDNS0_SUBNET)
	e.Code = dns.EDNS0SUBNET
	e.Family = uint16(family)
	e.SourceNetmask = uint8(bits)
	e.Address = addr.AsSlice()

	// OPT
	o := new(dns.OPT)
	o.Hdr.Name = "."
	o.Hdr.Rrtype = dns.TypeOPT
	o.Option = append(o.Option, e)

	return o, true
}

func containsEDNS0SUBNET(req *dns.Msg) bool {
	var found bool

	exdns.ForEachRR(req.Extra, func(opts *dns.OPT) {
		for _, opt := range opts.Option {
			if opt.Option() == dns.EDNS0SUBNET {
				found = true
			}
		}
	})

	return found
}

func removeEDNS0SUBNET(req *dns.Msg) bool {
	var altered bool

	filterEDNS0 := func(_ []dns.EDNS0, e dns.EDNS0) (dns.EDNS0, bool) {
		if e.Option() == dns.EDNS0SUBNET {
			// discard
			altered = true
			return nil, false
		}

		// keep
		return e, true
	}

	filterRR := func(_ []dns.RR, rr dns.RR) (dns.RR, bool) {
		if opts, ok := rr.(*dns.OPT); ok {
			// remove EDNS0SUBNET options
			opts.Option = core.SliceReplaceFn(opts.Option, filterEDNS0)

			if len(opts.Option) == 0 {
				// empty RR, discard
				altered = true
				return nil, false
			}
		}

		// keep
		return rr, true
	}

	req.Extra = core.SliceReplaceFn(req.Extra, filterRR)
	return altered
}
