package main

import (
	"fmt"
	"os"
	"path"

	"github.com/playlist-media/clu/config"
)

func getConfig() *config.Config {
	if cfg != nil {
		return cfg
	}

	// Get current directory
	pwd, err := os.Getwd()
	if err != nil {
		printFatal(fmt.Sprintf("Error: could not determine working directory (%s)", err.Error()))
	}

	// Parse config file
	configLocation := path.Join(pwd, "config", "clu.yaml")
	cfg, err = config.NewConfigFromFile(configLocation)
	if err != nil {
		printFatal(fmt.Sprintf("Error: could not parse the config file (./config/clu.yaml) (%s)\n", err.Error()))
	}

	return cfg
}
