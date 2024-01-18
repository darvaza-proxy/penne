// Package horizon provides the implementation for Penne Horizons
package horizon

import (
	"fmt"
	"net/http"
	"strings"

	"darvaza.org/core"
	"darvaza.org/sidecar/pkg/sidecar/horizon"

	"darvaza.org/penne/pkg/resolver"
)

// MakeHorizons builds Horizons from a [Config] slice,
// and prepares the resolvers to get back to us when
// they don't know what else to do.
func MakeHorizons(conf []Config,
	res map[string]*resolver.Resolver,
	ctxKey *core.ContextKey[horizon.Match]) ([]string, map[string]*Horizon, error) {
	//
	names, m, err := makeHorizonsMap(conf)
	if err != nil {
		return nil, nil, err
	}

	out := make(map[string]*Horizon)
	for len(m) > 0 {
		err := makeHorizonsPass(out, m, res, ctxKey)
		if err != nil {
			return nil, nil, err
		}
	}

	// hook fallback on all resolvers
	// so they can find their way back into the
	// horizons chain.
	for _, r := range res {
		r.SetFallback(dnsNextExchanger)
	}

	return names, out, nil
}

func makeHorizonsPass(out map[string]*Horizon, conf map[string]Config,
	res map[string]*resolver.Resolver, ctxKey *core.ContextKey[horizon.Match]) error {
	//
	name, next, err := nextMakeHorizons(out, conf)
	if err != nil {
		// broken dependencies tree
		return err
	}

	// take the chosen config and remove it from the map
	hc, ok := conf[name]
	if !ok {
		core.Panic("unreachable")
	}
	delete(conf, name)

	r, ok := getMakeHorizonsResolver(hc.Resolver, res)
	if !ok {
		// invalid resolver
		return &Error{
			Horizon: name,
			Reason:  fmt.Sprintf("resolver %q not found", hc.Resolver),
			Err:     core.ErrNotExists,
		}
	}

	h, err := hc.New(next, r, ctxKey)
	if err != nil {
		// failed to build horizon
		return err
	}

	// store
	out[name] = h
	return nil
}

func makeHorizonsMap(conf []Config) ([]string, map[string]Config, error) {
	names := make([]string, 0, len(conf))
	out := make(map[string]Config)

	for _, hc := range conf {
		hc.Name = strings.ToLower(hc.Name)
		hc.Next = strings.ToLower(hc.Next)
		hc.Resolver = strings.ToLower(hc.Resolver)

		if hc.Name == "" {
			err := &Error{
				Horizon: hc.Name,
				Reason:  "no name",
				Err:     core.ErrInvalid,
			}
			return nil, nil, err
		}

		if _, ok := out[hc.Name]; ok {
			err := &Error{
				Horizon: hc.Name,
				Reason:  "duplicate name",
				Err:     core.ErrExists,
			}
			return nil, nil, err
		}

		names = append(names, hc.Name)
		out[hc.Name] = hc
	}

	return names, out, nil
}

func nextMakeHorizons(out map[string]*Horizon,
	conf map[string]Config) (string, *Horizon, error) {
	//
	var err error

	for name, hc := range conf {
		// is the dependency is ready?
		z, ok := getMakeHorizons(hc.Next, out)
		if ok {
			// ready
			return name, z, nil
		}

		if err == nil {
			// first unresolvable
			err = &Error{
				Horizon: name,
				Reason:  fmt.Sprintf("horizon %q not found", hc.Next),
				Err:     core.ErrNotExists,
			}
		}
	}

	// none ready
	return "", nil, err
}

func getMakeHorizons(name string, out map[string]*Horizon) (*Horizon, bool) {
	//
	if name == "" {
		// no dependencies
		return nil, true
	}

	z, ok := out[name]
	return z, ok
}

func getMakeHorizonsResolver(name string, res map[string]*resolver.Resolver) (*resolver.Resolver, bool) {
	//
	if name == "" {
		// no resolver needed
		return nil, true
	}

	r, ok := res[name]
	return r, ok
}

// A Horizon is a [resolver.Exchanger] for a particular set of networks
type Horizon struct {
	next   *Horizon
	ctxKey *core.ContextKey[horizon.Match]
	res    *resolver.Resolver
	zc     horizon.Config

	allowForwarding bool

	nextH http.Handler
	nextE resolver.Exchanger
}

// New creates an assembled [horizon.Horizon]
func (z *Horizon) New(h http.Handler, e resolver.Exchanger) *horizon.Horizon {
	return z.zc.New(h, e)
}
