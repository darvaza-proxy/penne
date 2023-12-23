// Package resolver provides the implementation of Penne resolvers
package resolver

import (
	"darvaza.org/resolver"
)

var (
	_ resolver.Exchanger = (*Resolver)(nil)
)

// Resolver is a custom [resolver.Exchanger]
type Resolver struct {
	resolver.Exchanger
}
