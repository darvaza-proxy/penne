package suffix

import (
	"testing"

	"github.com/miekg/dns"
)

type suffixCase struct {
	s  string
	ok bool
	t  []suffixMatchCase
}

type suffixMatchCase struct {
	s  string
	ok bool
}

func tM(s string, ok bool) suffixMatchCase {
	return suffixMatchCase{
		s:  s,
		ok: ok,
	}
}

func tS(s string, ok bool, t ...suffixMatchCase) suffixCase {
	return suffixCase{
		s:  s,
		ok: ok,
		t:  t,
	}
}

func TestMatch(t *testing.T) {
	var cases = []suffixCase{
		tS("local", true,
			tM("local", true),
			tM("foo.local", true),
			tM("com", false),
			tM("example.com", false),
		),
	}
	for _, tc := range cases {
		testSuffixCase(t, tc)
	}
}

func testSuffixCase(t *testing.T, tc suffixCase) {
	p, err := Compile(tc.s)
	switch {
	case tc.ok && err != nil:
		t.Errorf("%q: failed unexpectedly: %v", tc.s, err)
	case !tc.ok && err != nil:
		t.Errorf("%q: failed as expected: %v", tc.s, err)
	case !tc.ok && err == nil:
		t.Errorf("%q: failed to fail", tc.s)
	default:
		testSuffixMatchCase(t, tc.s, &p, tc.t)
	}
}

func testSuffixMatchCase(t *testing.T, suffix string, p *Suffix, mc []suffixMatchCase) {
	for _, mcc := range mc {
		ok := p.Match(dns.CanonicalName(mcc.s))
		switch {
		case ok && !mcc.ok:
			t.Errorf("%q: %q: shouldn't have matched",
				suffix, mcc.s)
		case !ok && mcc.ok:
			t.Errorf("%q: %q: should have matched",
				suffix, mcc.s)
		case !ok:
			t.Logf("%q: %q: not a match",
				suffix, mcc.s)
		default:
			t.Logf("%q: %q: match",
				suffix, mcc.s)
		}
	}
}
