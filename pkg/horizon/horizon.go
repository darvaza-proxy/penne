// Package horizon provides the implementation for Penne Horizons
package horizon

import (
	"strings"

	"darvaza.org/core"
	"darvaza.org/resolver"
)

var (
	_ resolver.Exchanger = (*Horizon)(nil)
)

// MakeHorizons builds Horizons from a [Config] slice
func MakeHorizons(conf []Config,
	_ map[string]resolver.Exchanger) ([]string, map[string]*Horizon, error) {
	//
	names, _, err := makeHorizonsMap(conf)
	if err != nil {
		return nil, nil, err
	}

	out := make(map[string]*Horizon)
	return names, out, nil
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

// A Horizon is a [resolver.Exchanger] for a particular set of networks
type Horizon struct {
	resolver.Exchanger
}
