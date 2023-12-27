package main

import (
	"io"

	"github.com/spf13/pflag"

	"darvaza.org/core"
	"darvaza.org/sidecar/pkg/logger/zerolog"
	"darvaza.org/slog"
)

func newLogger(w io.Writer, flags *pflag.FlagSet) (slog.Logger, error) {
	level := slog.Error

	if flags != nil {
		verbosity, err := flags.GetCount(verbosityFlag)
		if err != nil {
			return nil, err
		}

		level += slog.LogLevel(verbosity)
		switch {
		case level < slog.Error:
			level = slog.Error
		case level > slog.Debug:
			level = slog.Debug
		}
	}

	log := zerolog.New(w, level)
	return log, nil
}

func mustLogger(w io.Writer, flags *pflag.FlagSet) slog.Logger {
	log, err := newLogger(w, flags)
	if err != nil {
		core.Panic(err)
	}
	return log
}

const (
	verbosityFlag      = "verbose"
	verbosityShortFlag = "v"
)

func init() {
	pFlags := rootCmd.PersistentFlags()
	pFlags.CountP(verbosityFlag, verbosityShortFlag, "increase verbosity")
}
