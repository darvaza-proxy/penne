package horizon

import (
	"net/http"
	"net/netip"

	"darvaza.org/core"
	"darvaza.org/sidecar/pkg/sidecar/horizon"

	"darvaza.org/penne/pkg/resolver"
)

// Config describes a [Horizon]
type Config struct {
	Name string `yaml:"name"`
	Next string `yaml:"next,omitempty" toml:",omitempty" json:",omitempty"`

	AllowForwarding bool `yaml:"allow_forwarding,omitempty" toml:",omitempty" json:",omitempty"`

	Networks []netip.Prefix `yaml:"networks,omitempty" toml:",omitempty" json:",omitempty"`
	Resolver string         `yaml:"resolver,omitempty" toml:",omitempty" json:",omitempty"`
}

// New creates a new [Horizon] from the [Config]
func (hc Config) New(next *Horizon, res *resolver.Resolver,
	ctxKey *core.ContextKey[horizon.Match]) (*Horizon, error) {
	//
	z := &Horizon{
		next: next,
		res:  res,

		ctxKey:          ctxKey,
		allowForwarding: hc.AllowForwarding,
	}

	z.zc = horizon.Config{
		Name:   hc.Name,
		Ranges: hc.Networks,

		Middleware:         newHorizonMiddleware(z),
		ExchangeMiddleware: newHorizonExchangeMiddleware(z),
	}

	return z, nil
}

func newHorizonMiddleware(z *Horizon) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// store last resort HTTP Handler
		z.nextH = next

		return http.HandlerFunc(z.HorizonServeHTTP)
	}
}

func newHorizonExchangeMiddleware(z *Horizon) func(resolver.Exchanger) resolver.Exchanger {
	return func(next resolver.Exchanger) resolver.Exchanger {
		// store last resort Exchanger
		z.nextE = next

		return resolver.ExchangerFunc(z.HorizonExchange)
	}
}
