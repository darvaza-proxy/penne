// Package suffix implement domain name suffix matching
package suffix

import (
	"strings"

	"darvaza.org/sidecar/pkg/glob"
)

// Compile converts an suffix pattern into a
// an [Suffix].
func Compile(expr string) (Suffix, error) {
	switch {
	case expr == "":
		// any
		expr = "."
	case expr[len(expr)-1] != '.':
		// not canonical
		expr += "."
	}

	if glob.HasGlobRunes(expr) {
		if expr[0] == '.' {
			return compileGlobSuffix("**" + expr)
		}

		return compileGlobSuffix("{**.,}" + expr)
	}

	if expr[0] == '.' {
		// .foo.bar -> **.foo.bar
		return compileGlobSuffix("**" + expr)
	}

	// literal
	return Suffix{l: expr}, nil
}

func compileGlobSuffix(pattern string) (Suffix, error) {
	g, err := glob.Compile(pattern, '.')
	if err != nil {
		return Suffix{}, err
	}

	return Suffix{g: g}, nil
}

// Suffix if a domain, or domain pattern, that
// can be matched against subdomains.
type Suffix struct {
	g *glob.Glob
	l string
}

// Match checks if a qName matches the [Suffix]
func (s Suffix) Match(qName string) bool {
	switch {
	case s.g != nil:
		return s.g.Match(qName)
	case s.l == qName:
		return true
	default:
		before, ok := strings.CutSuffix(qName, s.l)
		if ok && len(before) > 0 {
			return before[len(before)-1] == '.'
		}
	}
	return false
}

// Suffixes is a set of [Suffix] that restrict
// a [Resolver]. If empty, any name is accepted.
type Suffixes []Suffix

// Match checks if a qName matches any of the
// [Suffix]es. If none, it will always match
// as it means the [Resolver] is unrestricted.
func (ss Suffixes) Match(qName string) bool {
	if len(ss) == 0 {
		// unrestricted
		return true
	}

	for _, s := range ss {
		if s.Match(qName) {
			return true
		}
	}

	return false
}

// MakeSuffixes compiles a list of [Suffix]es.
func MakeSuffixes(suffixes []string) ([]Suffix, error) {
	out := make([]Suffix, 0, len(suffixes))
	for _, str := range suffixes {
		s, err := Compile(str)
		if err != nil {
			return nil, err
		}

		out = append(out, s)
	}
	return out, nil
}
