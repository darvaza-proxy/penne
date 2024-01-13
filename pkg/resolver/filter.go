package resolver

import (
	"context"

	"github.com/miekg/dns"

	"darvaza.org/core"
	"darvaza.org/resolver"
	"darvaza.org/resolver/pkg/exdns"
)

func newOmitEDNS0SubNetExchanger(next resolver.Exchanger) resolver.Exchanger {
	fn := func(ctx context.Context, req *dns.Msg) (*dns.Msg, error) {
		return omitEDNS0SubNetExchange(ctx, req, next)
	}

	return resolver.ExchangerFunc(fn)
}

func omitEDNS0SubNetExchange(ctx context.Context, req *dns.Msg, next resolver.Exchanger) (*dns.Msg, error) {
	var original = req

	req2 := req.Copy()
	if removeEDNS0SUBNET(req2) {
		req2.Id = dns.Id()
		req = req2
	}

	resp, err := next.Exchange(ctx, req)
	return exdns.RestoreReturn(original, resp, err)
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
