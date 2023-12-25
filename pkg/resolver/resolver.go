// Package resolver provides the implementation of Penne resolvers
package resolver

import (
	"fmt"
	"strings"

	"darvaza.org/core"
	"darvaza.org/resolver"
	"darvaza.org/slog"
)

var (
	_ resolver.Exchanger = (*Resolver)(nil)
)

// MakeResolvers builds resolvers from a [Config] slice.
func MakeResolvers(conf []Config,
	opts *Options) ([]string, map[string]resolver.Exchanger, error) {
	//
	if conf == nil {
		return nil, nil, core.ErrInvalid
	}

	names, m, err := makeResolversMap(conf)
	if err != nil {
		return nil, nil, err
	}

	res := make(map[string]resolver.Exchanger)
	for len(m) > 0 {
		err := makeResolversPass(res, m, opts)
		if err != nil {
			return nil, nil, err
		}
	}

	return names, res, nil
}

func makeResolversPass(res map[string]resolver.Exchanger,
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
			err := &Error{
				Resolver: rc.Name,
				Reason:   "no name",
				Err:      core.ErrInvalid,
			}
			return nil, nil, err
		}

		if _, ok := out[rc.Name]; ok {
			err := &Error{
				Resolver: rc.Name,
				Reason:   "duplicate name",
				Err:      core.ErrExists,
			}
			return nil, nil, err
		}

		names = append(names, rc.Name)
		out[rc.Name] = rc
	}

	return names, out, nil
}

func nextMakeResolvers(res map[string]resolver.Exchanger,
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

func getMakeResolvers(name string,
	res map[string]resolver.Exchanger) (resolver.Exchanger, bool) {
	//
	if name == "" {
		// no dependencies
		return nil, true
	}

	r, ok := res[name]
	return r, ok
}

// Resolver is a custom [resolver.Exchanger].
type Resolver struct {
	resolver.Exchanger

	log      slog.Logger
	name     string
	suffixes []string

	Next resolver.Exchanger
}

// Name returns the name of the resolver.
func (r *Resolver) Name() string {
	return r.name
}

func (r *Resolver) String() string {
	return fmt.Sprintf("resolver[%q]", r.name)
}
