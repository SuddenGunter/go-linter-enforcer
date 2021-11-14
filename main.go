package main

import (
	"io/ioutil"
	"os"

	"github.com/SuddenGunter/go-linter-enforcer/config"
	"github.com/SuddenGunter/go-linter-enforcer/enforcer"
	"github.com/SuddenGunter/go-linter-enforcer/logger"
	"github.com/SuddenGunter/go-linter-enforcer/repository"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"go.uber.org/zap"
)

func main() {
	log := logger.Create()
	cfg := config.FromEnv(log)

	repos, err := repository.LoadListFromJSON(cfg.RepositoriesFile)
	if err != nil {
		log.With("error", err).Fatal("unable to parse config file")
	}

	enf := enforcer.Enforcer{
		ExpectedFile: readAll(cfg.ExpectedLinterConfig, log),
		GitAuth: &http.BasicAuth{
			Username: cfg.Git.Username,
			Password: cfg.Git.Password,
		},
	}

	for _, r := range repos {
		enf.EnforceRules(r)
	}
}

func readAll(linterConfig string, log *zap.SugaredLogger) []byte {
	file, err := os.Open(linterConfig)
	if err != nil {
		log.Fatalw("failed to open file", "file", linterConfig, "err", err)
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalw("failed to read file", "file", linterConfig, "err", err)
	}

	return data
}
