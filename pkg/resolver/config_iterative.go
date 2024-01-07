package resolver

import (
	"darvaza.org/resolver"
	"darvaza.org/resolver/pkg/reflect"
	"darvaza.org/slog"
)

const (
	// DefaultIteratorCacheSize indicates the number of records
	// the Iterator lookuper will cache.
	DefaultIteratorCacheSize = 1024
)

func (rc Config) setupIterative(r *Resolver, opts *Options) error {
	var e resolver.Exchanger

	// TODO: give rc.Servers to the IteratorLookuper

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

	il := resolver.NewIteratorLookuper(rc.Name, DefaultIteratorCacheSize, c)
	il.SetLogger(opts.Logger)
	if err := il.AddRootServers(); err != nil {
		return &Error{
			Resolver: rc.Name,
			Reason:   "failed to create iterative lookuper",
			Err:      err,
		}
	}
	e = il

	if rc.OmitSubNet {
		e = newOmitEDNS0SubNetExchanger(e)
	}

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
