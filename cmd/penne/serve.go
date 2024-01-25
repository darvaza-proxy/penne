package main

import (
	"github.com/spf13/cobra"

	"darvaza.org/penne/pkg/server"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run DNS server",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		srv, err := server.New(srvConf)
		if err != nil {
			return err
		}

		return srv.ListenAndServe()
	},
}
