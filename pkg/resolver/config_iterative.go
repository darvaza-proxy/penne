package resolver

import (
	"darvaza.org/resolver"
	"darvaza.org/resolver/pkg/client"
	"darvaza.org/resolver/pkg/reflect"
	"darvaza.org/slog"
)

const (
	// DefaultIteratorMaxRR indicates the number of records
	// the Iterator lookuper will cache.
	DefaultIteratorMaxRR = 1024
)

func (rc Config) setupIterative(r *Resolver, opts *Options) error {
	c, err := rc.newClient(opts)
	if err != nil {
		return rc.WrapError(err, "failed to create client")
	}

	e, err := rc.newIteratorLookuper(c, opts)
	if err != nil {
		return rc.WrapError(err, "failed to create iterative lookuper")
	}

	// TODO: add cache

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

func (rc Config) newIteratorLookuper(c client.Client, opts *Options) (resolver.Exchanger, error) {
	var err error
	il := resolver.NewIteratorLookuper(rc.Name, rc.IterativeMaxRR, c)
	il.SetLogger(opts.Logger)

	if len(rc.Servers) == 0 {
		err = il.AddRootServers()
	} else {
		// use rc.Suffixes to narrow iteration to specific domains.
		err = il.AddServer(".", 0, rc.Servers...)
		_ = il.SetPersistent(".")
	}

	if err != nil {
		return nil, err
	}
	return il, nil
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
