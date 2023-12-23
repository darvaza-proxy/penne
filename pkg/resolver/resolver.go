// Package resolver provides the implementation of Penne resolvers
package resolver

import (
	"strings"

	"darvaza.org/core"
	"darvaza.org/resolver"
	"darvaza.org/slog"
)

var (
	_ resolver.Exchanger = (*Resolver)(nil)
)

// MakeResolvers builds resolvers from a [Config] slice
func MakeResolvers(conf []Config,
	_ slog.Logger) ([]string, map[string]resolver.Exchanger, error) {
	//
	names, _, err := makeResolversMap(conf)
	if err != nil {
		return nil, nil, err
	}

	res := make(map[string]resolver.Exchanger)
	return names, res, nil
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

// Resolver is a custom [resolver.Exchanger]
type Resolver struct {
	resolver.Exchanger
}
