// Package horizon provides the implementation for Penne Horizons
package horizon

import (
	"darvaza.org/resolver"
)

var (
	_ resolver.Exchanger = (*Horizon)(nil)
)

// A Horizon is a [resolver.Exchanger] for a particular set of networks
type Horizon struct {
	resolver.Exchanger
}
