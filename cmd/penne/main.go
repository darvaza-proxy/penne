// Package main implements the Penne server
package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"darvaza.org/sidecar/pkg/service"
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
		log.Print(err)
	}

	os.Exit(code)
}
