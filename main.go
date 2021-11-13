package main

import (
	"github.com/SuddenGunter/go-linter-enforcer/config"
	"github.com/go-git/go-git/v5/plumbing/transport/http"

	"github.com/SuddenGunter/go-linter-enforcer/repository"

	"github.com/SuddenGunter/go-linter-enforcer/logger"
)

func main() {
	log := logger.Create()
	cfg := config.FromEnv(log)

	repos, err := repository.LoadListFromJSON(cfg.RepositoriesFile)
	if err != nil {
		log.With("error", err).Fatal("unable to parse config file")
	}

	for _, r := range repos {
		err := repository.PushDemoBranch(&http.BasicAuth{
			Username: cfg.Git.Username,
			Password: cfg.Git.Password,
		}, r)
		if err != nil {
			log.With("error", err).Fatal("failed to push demo branch")
		}
	}
}
