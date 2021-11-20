package bitbucket

import (
	"io/ioutil"
	"os"
	"reflect"

	"github.com/SuddenGunter/go-linter-enforcer/git"
	"github.com/SuddenGunter/go-linter-enforcer/runner"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"go.uber.org/zap"
)

type RunnerBuilder struct{}

func (r RunnerBuilder) CreateRunner(log *zap.SugaredLogger, config interface{}) runner.Runner {
	cfg, ok := config.(Config)
	if !ok {
		log.Fatal("unable to assert config as bitbucket.Config{}")
	}

	gcp := git.NewClientProvider(log, &http.BasicAuth{
		Username: cfg.Git.Username,
		Password: cfg.Git.Password,
	})

	return NewRunner(gcp, readAll(cfg.ExpectedLinterConfig, log), log, cfg)
}

func (r RunnerBuilder) ConfigType() reflect.Type {
	return reflect.TypeOf(Config{})
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
