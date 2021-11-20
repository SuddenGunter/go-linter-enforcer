package raw

import (
	"io/ioutil"
	"os"

	"github.com/SuddenGunter/go-linter-enforcer/git"
	"github.com/SuddenGunter/go-linter-enforcer/repository"
	"github.com/SuddenGunter/go-linter-enforcer/runner"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"go.uber.org/zap"
)

type RunnerBuilder struct{}

func (r RunnerBuilder) CreateRunner(log *zap.SugaredLogger, config interface{}) runner.Runner {
	cfg, ok := config.(*Config)
	if !ok {
		log.Fatal("unable to assert config as raw.Config{}")
	}

	repos, err := repository.LoadListFromJSON(cfg.RepositoriesFile)
	if err != nil {
		log.With("error", err).Fatal("unable to parse repositories list file")
	}

	gcp := git.NewClientProvider(log, &http.BasicAuth{
		Username: cfg.Git.Username,
		Password: cfg.Git.Password,
	})

	return NewRunner(gcp, repos, readAll(cfg.ExpectedLinterConfig, log), log, *cfg)
}

func (r RunnerBuilder) Config() interface{} {
	return &Config{}
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
