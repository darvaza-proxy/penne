package main

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"darvaza.org/sidecar/pkg/glob"
	"darvaza.org/sidecar/pkg/service"
)

var globCmd = &cobra.Command{
	Use:   "glob",
	Short: "glob tests glob and replace patterns",

	RunE: func(cmd *cobra.Command, args []string) error {
		flags := cmd.Flags()

		// pattern
		g, args, err := globNewPattern(flags, args)
		if err != nil {
			return err
		}

		// replace
		r, err := globNewTemplate(flags)
		if err != nil {
			return err
		}

		// test
		for _, s := range args {
			globPrintResult(g, r, s)
		}

		return nil
	},
}

func globPrintResult(g *glob.Glob, t *glob.Template, fixture string) bool {
	data, ok := g.Capture(fixture)
	if !ok {
		_, _ = fmt.Printf("- %q %s\n", fixture, "FAIL")
		return false
	}

	_, _ = fmt.Printf("- %q %s\n", fixture, "MATCH")
	for i, s := range data {
		_, _ = fmt.Printf("  %v: %q\n", i+1, s)
	}

	if t != nil {
		s, err := t.Replace(data)
		if err != nil {
			_, _ = fmt.Printf("  => %v\n", err)
			return false
		}

		_, _ = fmt.Printf("  => %q\n", s)
	}

	return true
}

func globNewPattern(flags *pflag.FlagSet, args []string) (*glob.Glob, []string, error) {
	flag := flags.Lookup(globPatternFlag)
	switch {
	case flag == nil:
		panic("unreachable")
	case flag.Changed:
		s := flag.Value.String()
		_, _ = fmt.Println("Glob:", s)
		g, err := glob.Compile(s, '.')
		return g, args, err
	case len(args) > 0:
		s := args[0]
		args = args[1:]
		_, _ = fmt.Println("Glob:", s)
		g, err := glob.Compile(s, '.')
		return g, args, err
	default:
		err := &service.ErrorExitCode{
			Code: 1,
			Err:  errors.New("no pattern specified"),
		}
		return nil, nil, err
	}
}

func globNewTemplate(flags *pflag.FlagSet) (*glob.Template, error) {
	flag := flags.Lookup(globTemplateFlag)
	switch {
	case flag == nil:
		panic("unreachable")
	case flag.Changed:
		s := flag.Value.String()
		_, _ = fmt.Println("Template:", s)
		return glob.CompileTemplate(s)
	default:
		return nil, nil
	}
}

const (
	globPatternFlag       = "glob"
	globPatternShortFlag  = "g"
	globTemplateFlag      = "replace"
	globTemplateShortFlag = "r"
)

func init() {
	flags := globCmd.Flags()
	flags.StringP(globPatternFlag, globPatternShortFlag, "",
		"glob pattern to test")
	flags.StringP(globTemplateFlag, globTemplateShortFlag, "",
		"template where to apply glob captures")

	rootCmd.AddCommand(globCmd)
}
