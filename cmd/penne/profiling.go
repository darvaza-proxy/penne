package main

import (
	"os"
	"runtime/pprof"

	"darvaza.org/slog"
	"github.com/spf13/cobra"
)

func memProfilingInit() {
	flags := rootCmd.Flags()

	mpf := flags.Lookup(memProfileFlag)
	if mpf.Changed {
		log := newLogger(flags)
		memProfile := mpf.Value.String()

		log.Info().
			WithField("fileName", memProfile).
			Print("memory profile enabled")

		// Create profiling file
		f, err := os.Create(memProfile)
		if err != nil {
			log.Error().
				WithField(slog.ErrorFieldName, err).
				WithField("fileName", memProfile).
				Print("failed to create memory profile file")
			return
		}

		// Schedule writing of memory profile
		cobra.OnFinalize(func() {
			defer f.Close()
			_ = pprof.WriteHeapProfile(f)

			log.Info().
				WithField("fileName", memProfile).
				Print("memory profile stopped")
		})
	}
}

func cpuProfilingInit() {
	flags := rootCmd.Flags()

	cpf := flags.Lookup(cpuProfileFlag)
	if cpf.Changed {
		log := newLogger(flags)
		cpuProfile := cpf.Value.String()

		log.Info().
			WithField("fileName", cpuProfile).
			Print("CPU profile enabled")

		// Create profiling file
		f, err := os.Create(cpuProfile)
		if err != nil {
			log.Error().
				WithField(slog.ErrorFieldName, err).
				WithField("fileName", cpuProfile).
				Print("failed to create CPU profile file")
			return
		}

		// Start profiling
		err = pprof.StartCPUProfile(f)
		if err != nil {
			_ = f.Close()

			log.Error().
				WithField(slog.ErrorFieldName, err).
				WithField("fileName", cpuProfile).
				Print("failed to start CPU profiling")
			return
		}

		// Schedule stop
		cobra.OnFinalize(func() {
			pprof.StopCPUProfile()
			_ = f.Close()

			log.Info().
				WithField("fileName", cpuProfile).
				Print("CPU profile stopped")
		})
	}
}

const (
	cpuProfileFlag = "cpu-profile"
	memProfileFlag = "mem-profile"
)

func init() {
	pFlags := rootCmd.PersistentFlags()
	pFlags.String(cpuProfileFlag, "", "write CPU profile to file")
	pFlags.String(memProfileFlag, "", "write MEM profile to file")

	cobra.OnInitialize(cpuProfilingInit)
	cobra.OnInitialize(memProfilingInit)
}
