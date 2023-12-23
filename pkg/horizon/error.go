package horizon

import (
	"fmt"
	"strings"

	"darvaza.org/core"
)

var (
	_ core.Unwrappable = (*Error)(nil)
)

// Error is an error that references the name of a Horizon
type Error struct {
	Horizon string
	Reason  string
	Err     error
}

func (e Error) Error() string {
	s := make([]string, 0, 3)

	s = append(s, fmt.Sprintf("horizon[%q]", e.Horizon))
	if e.Reason != "" {
		s = append(s, e.Reason)
	}
	if e.Err != nil {
		s = append(s, e.Err.Error())
	}

	return strings.Join(s, ": ")
}

func (e Error) Unwrap() error {
	return e.Err
}
