package resolver

import (
	"darvaza.org/resolver"
	"darvaza.org/resolver/pkg/reflect"
	"darvaza.org/slog"
)

func (rc Config) setupIterative(r *Resolver, opts *Options) error {
	var e resolver.Exchanger

	if rc.Recursive || len(rc.Servers) > 0 {
		return &Error{
			Resolver: rc.Name,
			Reason:   "iterative resolver with specific servers or recursive not supported.",
		}
	}

	c, err := rc.newClient(opts)
	if err != nil {
		return &Error{
			Resolver: rc.Name,
			Reason:   "failed to create client",
			Err:      err,
		}
	}

	e, err = resolver.NewRootLookuperWithClient("", c)
	if err != nil {
		return &Error{
			Resolver: rc.Name,
			Reason:   "failed to create iterative lookuper",
			Err:      err,
		}
	}

	// TODO: add cache

	if rc.Debug {
		e, _ = reflect.NewWithExchanger(rc.Name, opts.Logger, e)

		rc.setupIterativeDebug(r)
	}

	r.Exchanger = e
	return nil
}

func (rc Config) setupIterativeDebug(r *Resolver) {
	// Info level for remote calls
	r.debug[rc.Name+"-udp"] = slog.Info
	r.debug[rc.Name+"-tcp"] = slog.Info
	r.debug[rc.Name+"-tls"] = slog.Info
	// Debug for those that could be cached
	r.debug[rc.Name+"-mux"] = slog.Debug
	r.debug[rc.Name] = slog.Debug
}
