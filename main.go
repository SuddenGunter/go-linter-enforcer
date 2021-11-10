package main

import (
	"os"

	"github.com/go-git/go-git/v5/plumbing/transport/http"

	"github.com/SuddenGunter/go-linter-enforcer/repository"

	"github.com/SuddenGunter/go-linter-enforcer/logger"
)

func main() {
	log := logger.Create()

	cfgFile := os.Getenv("CONFIG_FILE")
	if cfgFile == "" {
		log.Fatal("CONFIG_FILE environment variable is required")
	}

	demoPass := os.Getenv("DEMO_PASSWORD")
	if demoPass == "" {
		log.Fatal("DEMO_PASSWORD environment variable is required")
	}

	demoUser := os.Getenv("DEMO_USERNAME")
	if demoUser == "" {
		log.Fatal("DEMO_USERNAME environment variable is required")
	}

	cfg, err := repository.ConfigFromJSON(cfgFile)
	if err != nil {
		log.With("error", err).Fatal("unable to parse config file")
	}

	for _, r := range cfg.Repositories {
		err := repository.PushDemoBranch(&http.BasicAuth{
			Username: demoUser,
			Password: demoPass,
		}, r)
		if err != nil {
			log.With("error", err).Fatal("failed to push demo branch")
		}
	}
}
