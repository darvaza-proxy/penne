package main

import (
	"context"
	"io/fs"
	"path/filepath"

	"github.com/spf13/pflag"

	"darvaza.org/sidecar/pkg/config"
	"darvaza.org/slog"

	"darvaza.org/penne/pkg/server"
)

const (
	configFileFlag      = "config"
	configFileShortFlag = "f"
	configFileDefault   = CmdName + ".{conf,json,toml,yaml}"
)

var confLoader = config.Loader[server.Config]{
	Base:       CmdName,
	Extensions: []string{"conf", "json", "toml", "yaml"},
}

func getConfig(ctx context.Context, flags *pflag.FlagSet) (fs.FS, *server.Config, error) {
	log := getLogger(ctx, flags)
	init := func(cfg *server.Config) error {
		cfg.Context = ctx
		cfg.Logger = log
		return nil
	}

	flag := flags.Lookup(configFileFlag)
	cfg, err := confLoader.NewFromFlag(flag, init)
	if err != nil {
		log.Error().WithField(slog.ErrorFieldName, err).Print("LoadConfigFile")
		return nil, nil, err
	}

	if fSys, fileName := confLoader.Last(); fileName != "" {
		log.Info().WithField("filename", fileName).Print("config loaded")

		fSys, err = fs.Sub(fSys, filepath.Dir(fileName))
		if err == nil {
			return fSys, cfg, nil
		}
	}

	return nil, cfg, nil
}

func init() {
	pFlags := rootCmd.PersistentFlags()
	pFlags.StringP(configFileFlag, configFileShortFlag, configFileDefault, "config file to use")
}
