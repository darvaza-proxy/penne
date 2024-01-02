package main

import (
	"os"
	"runtime/pprof"

	"darvaza.org/slog"
)

func cpuProfilingInit() {
	flags := rootCmd.Flags()

	cpf := flags.Lookup(cpuProfileFlag)
	if cpf.Changed {
		log := mustLogger(nil, flags)
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
				Print("failed to create CPU profile fine")
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
		doOnFinalize(func() {
			pprof.StopCPUProfile()
			_ = f.Close()

			log.Info().
				WithField("fileName", cpuProfile).
				Print("CPU profile stopped")
		})
	}
}

const cpuProfileFlag = "cpu-profile"

func init() {
	pFlags := rootCmd.PersistentFlags()
	pFlags.String(cpuProfileFlag, "", "write CPU profile to file")

	doOnInit(cpuProfilingInit)
}
