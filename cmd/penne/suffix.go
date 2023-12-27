package main

import (
	"errors"
	"fmt"

	"github.com/miekg/dns"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"darvaza.org/penne/pkg/suffix"
	"darvaza.org/sidecar/pkg/service"
)

var suffixCmd = &cobra.Command{
	Use:   "suffix",
	Short: "suffix tests suffix patterns",

	RunE: func(cmd *cobra.Command, args []string) error {
		flags := cmd.Flags()

		// suffix
		p, args, err := suffixNewPattern(flags, args)
		if err != nil {
			return err
		}

		// test
		for _, s := range args {
			suffixPrintResult(p, s)
		}

		return nil
	},
}

func suffixPrintResult(p *suffix.Suffix, fixture string) bool {
	var s string

	ok := p.Match(dns.CanonicalName(fixture))
	if ok {
		s = "MATCH"
	} else {
		s = "FAIL"
	}

	_, _ = fmt.Printf("- %q %s\n", fixture, s)
	return ok
}

func suffixNewPattern(flags *pflag.FlagSet, args []string) (*suffix.Suffix, []string, error) {
	var s string
	var ok bool

	flag := flags.Lookup(suffixPatternFlag)

	switch {
	case flag == nil:
		panic("unreachable")
	case flag.Changed:
		s = flag.Value.String()
		ok = true
	case len(args) > 0:
		s = args[0]
		args = args[1:]
		ok = true
	}

	if !ok {
		err := &service.ErrorExitCode{
			Code: 1,
			Err:  errors.New("no suffix specified"),
		}
		return nil, nil, err
	}

	_, _ = fmt.Println("Suffix:", s)
	p, err := suffix.Compile(s)
	if err != nil {
		return nil, nil, err
	}

	return &p, args, nil
}

const (
	suffixPatternFlag      = "suffix"
	suffixPatternShortFlag = "s"
)

func init() {
	flags := suffixCmd.Flags()
	flags.StringP(suffixPatternFlag, suffixPatternShortFlag, "",
		"suffix to test")
	rootCmd.AddCommand(suffixCmd)
}
