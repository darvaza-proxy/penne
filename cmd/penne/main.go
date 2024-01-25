// Package main implements the Penne server
package main

import (
	"os"

	"github.com/spf13/cobra"

	"darvaza.org/sidecar/pkg/service"
	"darvaza.org/slog"

	"darvaza.org/penne/pkg/server"
)

const (
	// CmdName is the name of this executable
	CmdName = "penne"
)

var rootCmd = &cobra.Command{
	Use:   CmdName,
	Short: "penne resolves names",
	Args:  cobra.NoArgs,

	PersistentPreRunE: setup,

	SilenceErrors: true,
	SilenceUsage:  true,
}

var srvConf *server.Config

func setup(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	flags := cmd.Flags()

	cfg, err := getConfig(ctx, flags)
	if err != nil {
		return err
	}

	// store
	srvConf = cfg
	return nil
}

func main() {
	err := rootCmd.Execute()
	code, err := service.AsExitStatus(err)

	if err != nil {
		newLogger(nil).Error().
			WithField(slog.ErrorFieldName, err).
			Print()
	}

	os.Exit(code)
}
