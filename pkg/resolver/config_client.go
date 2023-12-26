package resolver

import (
	"time"

	"darvaza.org/resolver/pkg/client"
	"darvaza.org/resolver/pkg/reflect"
)

func (rc Config) newClient(opts *Options) (client.Client, error) {
	var c, udp, tcp client.Client
	var err error

	// UDP/TCP client mux
	udp = opts.NewClient("udp")
	tcp = opts.NewClient("tcp")
	if rc.Debug {
		udp, _ = reflect.NewWithClient(rc.Name+"-udp", opts.Logger, udp)
		tcp, _ = reflect.NewWithClient(rc.Name+"-tcp", opts.Logger, tcp)
	}

	c, err = client.NewAutoClient(udp, tcp, 1*time.Second)
	if err != nil {
		return nil, err
	}

	// TODO: add DNS over TLS support
	// TODO: add DNS over HTTPS support

	if rc.DisableAAAA {
		// remove AAAA entries
		c = client.NewNoAAAA(c)
	}

	if rc.Debug {
		c, _ = reflect.NewWithClient(rc.Name+"-mux", opts.Logger, c)
	}

	return c, nil
}
