package resolver

import (
	"context"

	"github.com/miekg/dns"

	"darvaza.org/resolver"
	"darvaza.org/resolver/pkg/client"
	"darvaza.org/resolver/pkg/exdns"
	"darvaza.org/resolver/pkg/reflect"
	"darvaza.org/slog"
)

func (rc Config) setupForwarder(r *Resolver, opts *Options) error {
	c, err := rc.newClient(opts)
	if err != nil {
		return rc.WrapError(err, "failed to create client")
	}

	e, err := rc.newForwardLookuper(c, opts)
	if err != nil {
		return rc.WrapError(err, "failed to create forward lookuper")
	}

	// TODO: add cache

	if rc.OmitSubNet {
		e = newOmitEDNS0SubNetExchanger(e)
	}

	if rc.Debug {
		e, _ = reflect.NewWithExchanger(rc.Name, opts.Logger, e)

		rc.setupForwardDebug(r)
	}

	r.e = e
	return nil
}

func (rc Config) newForwardLookuper(c client.Client, opts *Options) (resolver.Exchanger, error) {
	var e resolver.Exchanger
	e, err := resolver.NewPoolExchanger(c, rc.Servers...)
	if err != nil {
		return nil, err
	}

	if opts.SingleFlight > 0 {
		e, err = resolver.NewSingleFlight(e, opts.SingleFlight, nil)
		if err != nil {
			return nil, err
		}
	}

	e = newRecursionDesired(e, rc.Recursive)
	return e, nil
}

func (rc Config) setupForwardDebug(r *Resolver) {
	rc.setupClientDebug(r, slog.Info, slog.Debug)
}

// revive:disable:flag-parameter

func newRecursionDesired(next resolver.Exchanger, recursive bool) resolver.ExchangerFunc {
	// revive:enable:flag-parameter
	return func(ctx context.Context, req *dns.Msg) (*dns.Msg, error) {
		if req.RecursionDesired != recursive {
			// change RecursionDesired flag to match the Resolver Config.
			req2 := req.Copy()
			req2.Id = dns.Id()
			req2.RecursionDesired = recursive

			resp, err := next.Exchange(ctx, req2)
			return exdns.RestoreReturn(req, resp, err)
		}

		return next.Exchange(ctx, req)
	}
}
