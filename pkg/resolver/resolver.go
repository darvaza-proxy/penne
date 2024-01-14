// Package resolver provides the implementation of Penne resolvers
package resolver

import (
	"context"
	"fmt"
	"strings"

	"github.com/miekg/dns"

	"darvaza.org/core"
	"darvaza.org/resolver"
	"darvaza.org/resolver/pkg/client"
	"darvaza.org/resolver/pkg/errors"
	"darvaza.org/resolver/pkg/exdns"
	"darvaza.org/slog"

	"darvaza.org/penne/pkg/suffix"
)

var (
	_ resolver.Lookuper  = (*Resolver)(nil)
	_ resolver.Exchanger = (*Resolver)(nil)
)

// MakeResolvers builds resolvers from a [Config] slice.
func MakeResolvers(conf []Config, debug map[string]slog.LogLevel,
	opts *Options) ([]string, map[string]*Resolver, error) {
	//
	if conf == nil {
		return nil, nil, core.ErrInvalid
	}

	names, m, err := makeResolversMap(conf)
	if err != nil {
		return nil, nil, err
	}

	res := make(map[string]*Resolver)
	for len(m) > 0 {
		err := makeResolversPass(res, m, opts)
		if err != nil {
			return nil, nil, err
		}
	}

	if debug != nil {
		makeResolverSetDebug(res, debug)
	}
	return names, res, nil
}

func makeResolversPass(res map[string]*Resolver,
	conf map[string]Config, opt *Options) error {
	//
	name, next, err := nextMakeResolvers(res, conf)
	if err != nil {
		// broken dependencies tree
		return err
	}

	// take the chosen config and remove it from the map
	rc, ok := conf[name]
	if !ok {
		core.Panic("unreachable")
	}
	delete(conf, name)

	r, err := rc.New(next, opt)
	if err != nil {
		// failed to build resolver
		return err
	}

	// store
	res[name] = r
	return nil
}

func makeResolversMap(conf []Config) ([]string, map[string]Config, error) {
	names := make([]string, 0, len(conf))
	out := make(map[string]Config)

	for _, rc := range conf {
		rc.Name = strings.ToLower(rc.Name)
		rc.Next = strings.ToLower(rc.Next)

		if rc.Name == "" {
			err := rc.WrapError(core.ErrInvalid, "no name")
			return nil, nil, err
		}

		if _, ok := out[rc.Name]; ok {
			err := rc.WrapError(core.ErrExists, "duplicate name")
			return nil, nil, err
		}

		names = append(names, rc.Name)
		out[rc.Name] = rc
	}

	return names, out, nil
}

func makeResolverSetDebug(res map[string]*Resolver,
	debug map[string]slog.LogLevel) {
	//
	for _, r := range res {
		r.copyDebugMap(debug)
	}
}

func nextMakeResolvers(res map[string]*Resolver,
	conf map[string]Config) (string, resolver.Exchanger, error) {
	//
	var err error

	for name, rc := range conf {
		// is the dependency ready?
		r, ok := getMakeResolvers(rc.Next, res)
		if ok {
			// ready
			return name, r, nil
		}

		if err == nil {
			// first unresolvable
			err = &Error{
				Resolver: name,
				Reason:   fmt.Sprintf("resolver %q not available", rc.Next),
				Err:      core.ErrNotExists,
			}
		}
	}

	// none ready
	return "", nil, err
}

func getMakeResolvers(name string, res map[string]*Resolver) (*Resolver, bool) {
	if name == "" {
		// no dependencies
		return nil, true
	}

	r, ok := res[name]
	return r, ok
}

// Resolver is a custom [resolver.Exchanger].
type Resolver struct {
	debug    map[string]slog.LogLevel
	log      slog.Logger
	name     string
	rewrite  Rewriters
	suffixes suffix.Suffixes
	w        client.Worker
	e        resolver.Exchanger
	next     resolver.Exchanger
}

// Name returns the name of the resolver.
func (r *Resolver) Name() string {
	return r.name
}

func (r *Resolver) String() string {
	return fmt.Sprintf("resolver[%q]", r.name)
}

// Start starts the [Resolver]'s worker.
func (r *Resolver) Start(ctx context.Context) error {
	if r != nil && r.w != nil {
		return r.w.Start(ctx)
	}
	return nil
}

// Cancel initiates a shut down of [Resolver]'s Worker.
func (r *Resolver) Cancel(err error) bool {
	if r != nil && r.w != nil {
		return r.w.Cancel(err)
	}
	return false
}

// Shutdown initiates a shut down of [Resolver]'s Worker,
// and waits until they are done or the given context expires.
func (r *Resolver) Shutdown(ctx context.Context) error {
	if r != nil && r.w != nil {
		return r.w.Shutdown(ctx)
	}
	return nil
}

func (r *Resolver) copyDebugMap(debug map[string]slog.LogLevel) bool {
	if len(r.debug) > 0 {
		for k, v := range r.debug {
			debug[k] = v
		}
		return true
	}
	return false
}

// Lookup implements the [resolver.Lookuper] interface.
func (r *Resolver) Lookup(ctx context.Context, qName string, qType uint16) (*dns.Msg, error) {
	req := exdns.NewRequestFromParts(qName, dns.ClassINET, qType)
	return r.Exchange(ctx, req)
}

// Exchange implements the [resolver.Exchanger] interface.
func (r *Resolver) Exchange(ctx context.Context, req *dns.Msg) (*dns.Msg, error) {
	switch {
	case ctx == nil, req == nil, len(req.Question) != 1:
		// This is unreachable. The server won't pass
		// bad requests across.
		return nil, errors.ErrBadRequest()
	case !r.match(req):
		// carry on, nothing to see here
		return r.nextExchange(ctx, req)
	case len(r.rewrite) == 0:
		// use our exchanger if available
		return r.doExchange(ctx, req)
	default:
		// rewrites involved, time to work
		return r.rewriteExchange(ctx, req)
	}
}

func (*Resolver) rewriteExchange(context.Context, *dns.Msg) (*dns.Msg, error) {
	// TODO: apply rewrites to questions
	// TODO: restore questions
	// TODO: apply rewrites to answers, and ask again if needed.
	return nil, errors.ErrNotImplemented("")
}

func (r *Resolver) doExchange(ctx context.Context, req *dns.Msg) (*dns.Msg, error) {
	if r.e != nil {
		return r.e.Exchange(ctx, req)
	}

	return r.nextExchange(ctx, req)
}

func (r *Resolver) nextExchange(ctx context.Context, req *dns.Msg) (*dns.Msg, error) {
	if r.next != nil {
		return r.next.Exchange(ctx, req)
	}

	return forbiddenExchange(ctx, req)
}

func (r *Resolver) match(req *dns.Msg) bool {
	if len(r.suffixes) == 0 {
		return true
	}

	q := req.Question[0]
	return r.suffixes.Match(q.Name)
}

// SetFallback sets the exchanger to use next if it doesn't have
// one already set from [Config].
func (r *Resolver) SetFallback(last resolver.Exchanger) bool {
	if r.next == nil && last != nil {
		r.next = last
		return true
	}

	return false
}

func forbiddenExchange(_ context.Context, req *dns.Msg) (*dns.Msg, error) {
	resp := new(dns.Msg)
	resp.SetRcode(req, dns.RcodeRefused)
	resp.Compress = false
	resp.RecursionAvailable = true
	return resp, nil
}
