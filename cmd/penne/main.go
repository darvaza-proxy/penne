// Package main implements the Penne server
package main

import (
	"os"

	"github.com/spf13/cobra"

	"darvaza.org/sidecar/pkg/service"
	"darvaza.org/slog"
)

const (
	// CmdName is the name of this executable
	CmdName = "penne"
)

var rootCmd = &cobra.Command{
	Use:   CmdName,
	Short: "penne resolves names",
}

func main() {
	err := rootCmd.Execute()
	code, err := service.AsExitStatus(err)

	if err != nil {
		mustLogger(nil, nil).Error().
			WithField(slog.ErrorFieldName, err).
			Print()
	}

	os.Exit(code)
}
