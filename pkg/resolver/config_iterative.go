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
	w, c, err := rc.newClient(opts)
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

	r.w = w
	r.e = e
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
	rc.setupClientDebug(r, slog.Info, slog.Debug)
	r.debug[rc.Name] = slog.Debug
}
