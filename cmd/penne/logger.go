package main

import (
	"github.com/spf13/pflag"

	"darvaza.org/sidecar/pkg/logger/zerolog"
	"darvaza.org/slog"
)

func getLogLevel(flags *pflag.FlagSet) slog.LogLevel {
	level := slog.Error

	if flags != nil {
		verbosity, err := flags.GetCount(verbosityFlag)
		if err == nil {
			level += slog.LogLevel(verbosity)
			switch {
			case level < slog.Error:
				level = slog.Error
			case level > slog.Debug:
				level = slog.Debug
			}
		}
	}

	return level
}

func newLogger(flags *pflag.FlagSet) slog.Logger {
	return newLoggerLevel(getLogLevel(flags))
}

func newLoggerLevel(level slog.LogLevel) slog.Logger {
	return zerolog.New(nil, level)
}

const (
	verbosityFlag      = "verbose"
	verbosityShortFlag = "v"
)

func init() {
	pFlags := rootCmd.PersistentFlags()
	pFlags.CountP(verbosityFlag, verbosityShortFlag, "increase verbosity")
}
