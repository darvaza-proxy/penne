package resolver

import (
	"time"

	"darvaza.org/resolver/pkg/client"
	"darvaza.org/resolver/pkg/reflect"
)

func (rc Config) newClient(opts *Options) (client.Client, error) {
	var c, udp, tcp client.Client

	// UDP/TCP client mux
	udp = opts.NewClient("udp")
	tcp = opts.NewClient("tcp")
	if rc.Debug {
		udp, _ = reflect.NewWithClient(rc.Name+"-udp", opts.Logger, udp)
		tcp, _ = reflect.NewWithClient(rc.Name+"-tcp", opts.Logger, tcp)
	}

	mux, err := client.NewAutoClient(udp, tcp, 1*time.Second)
	if err != nil {
		return nil, err
	}

	// DNS over TLS
	c = opts.NewClient("tcp+tls")
	switch {
	case c == nil:
		// not TLS
	case rc.Debug:
		// reflect TLS
		mux.TLS, _ = reflect.NewWithClient(rc.Name+"-tls", opts.Logger, c)
	default:
		// direct TLS
		mux.TLS = c
	}

	// TODO: add DNS over HTTPS support

	c = mux
	if rc.DisableAAAA {
		// remove AAAA entries
		c = client.NewNoAAAA(c)
	}

	if opts.SingleFlight >= 0 {
		// stampede control
		c = client.NewSingleFlight(c, opts.SingleFlight)
	}

	if rc.Debug {
		// logging
		c, _ = reflect.NewWithClient(rc.Name+"-mux", opts.Logger, c)
	}

	return c, nil
}
