package main

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"darvaza.org/penne/pkg/server"
	"darvaza.org/sidecar/pkg/service"
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

func setupService(_ context.Context, _ *service.Service, _ *server.Config) error {
	// TODO: implement
	return nil
}

// WantsSyslog tells if the `--syslog` flag was passed
// to use the system logger in interactive mode.
func WantsSyslog(flags *pflag.FlagSet) bool {
	v, _ := flags.GetBool(syslogFlag)
	return v
}

const syslogFlag = "syslog"

func init() {
	flags := serveCmd.Flags()
	flags.Bool(syslogFlag, false, "use syslog when running manually")
}
