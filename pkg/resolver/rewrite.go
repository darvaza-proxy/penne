package resolver

import "darvaza.org/core"

// Rewriter ...
type Rewriter struct{}

// Rewriters ...
type Rewriters []Rewriter

// MakeRewriters ...
func MakeRewriters(rcc []RewriteConfig) ([]Rewriter, error) {
	out := make([]Rewriter, len(rcc))
	for range rcc {
		return nil, core.ErrNotImplemented
	}
	return out, nil
}
