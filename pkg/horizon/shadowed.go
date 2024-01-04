package horizon

import (
	"darvaza.org/core"
	"darvaza.org/sidecar/pkg/sidecar/horizon"
)

// NewContextKey returns a [core.ContextKey] to be used
// to store the [horizon.Match] in [horizon.Horizons]
func NewContextKey(name string) *core.ContextKey[horizon.Match] {
	// pass-through in the case of shadowing
	return horizon.NewContextKey(name)
}
