package main

import (
	"github.com/SuddenGunter/go-linter-enforcer/config"
	"github.com/SuddenGunter/go-linter-enforcer/enforcer"
	"github.com/SuddenGunter/go-linter-enforcer/logger"
	"github.com/SuddenGunter/go-linter-enforcer/repository"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func main() {
	log := logger.Create()
	cfg := config.FromEnv(log)

	repos, err := repository.LoadListFromJSON(cfg.RepositoriesFile)
	if err != nil {
		log.With("error", err).Fatal("unable to parse config file")
	}

	enf := enforcer.Enforcer{
		GitAuth: &http.BasicAuth{
			Username: cfg.Git.Username,
			Password: cfg.Git.Password,
		},
	}

	for _, r := range repos {
		if err := enf.EnforceRules(r); err != nil {
			log.With("error", err).Fatal("failed to push demo branch")
		}
	}
}
