package main

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"darvaza.org/sidecar/pkg/config"
)

const (
	// DefaultDumpFormat indicates the format used by `penne dump` if none
	// is specified.
	DefaultDumpFormat = "toml"

	dumpFormatFlag      = "format"
	dumpFormatShortFlag = "T"
)

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "prints out the configuration",
	Args:  cobra.NoArgs,

	RunE: func(cmd *cobra.Command, _ []string) error {
		flags := cmd.Flags()

		ctx := context.TODO()
		cfg, err := prepareConfig(ctx, flags)
		if err != nil {
			return err
		}

		encFormat, err := flags.GetString(dumpFormatFlag)
		if err != nil {
			return err
		}

		enc, err := config.NewEncoder(encFormat)
		if err != nil {
			return err
		}

		_, err = enc.WriteTo(cfg, os.Stdout)
		return err
	},
}

func init() {
	flags := dumpCmd.Flags()
	flags.StringP(dumpFormatFlag, dumpFormatShortFlag, DefaultDumpFormat,
		"file format for the dump (yaml, toml or json)")

	rootCmd.AddCommand(dumpCmd)
}
