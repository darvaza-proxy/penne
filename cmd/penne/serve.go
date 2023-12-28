package main

import (
	"context"

	"github.com/spf13/cobra"

	"darvaza.org/penne/pkg/server"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run DNS server",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := context.Background()
		cfg, err := prepareConfig(ctx, cmd.Flags())
		if err != nil {
			return err
		}

		srv, err := server.New(cfg)
		if err != nil {
			return err
		}

		return srv.ListenAndServe()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
