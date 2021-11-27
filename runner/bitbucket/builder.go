package bitbucket

import (
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/go-git/go-git/v5/plumbing/transport/ssh"

	"github.com/SuddenGunter/go-linter-enforcer/git"
	"github.com/SuddenGunter/go-linter-enforcer/runner"
	"go.uber.org/zap"
)

type RunnerBuilder struct{}

func (r RunnerBuilder) CreateRunner(log *zap.SugaredLogger, config interface{}) runner.Runner {
	cfg, ok := config.(*Config)
	if !ok {
		log.Fatal("unable to assert config as bitbucket.Config{}")
	}

	publicKeys, err := ssh.NewPublicKeysFromFile("git", cfg.Git.SSHPrivateKeyPath, cfg.Git.SSHPrivateKeyPassword)
	if err != nil {
		log.Fatalw("failed to configure ssh", "SSHPrivateKeyPath", cfg.Git.SSHPrivateKeyPath, "err", err)
	}

	gcp := git.NewClientProvider(log, publicKeys)

	apiClient := getAPIClient(cfg, log)

	return NewRunner(gcp, readAll(cfg.ExpectedLinterConfig, log), log, apiClient, cfg)
}

func getAPIClient(cfg *Config, log *zap.SugaredLogger) APIClient {
	client := &Client{
		Client: http.Client{
			Timeout: 15 * time.Second,
		},
		Organization: cfg.Organization,
		Login:        cfg.Login,
		AppPassword:  cfg.AppPassword,
	}

	if cfg.DryRun {
		return UseDryRun(client, log)
	}

	return client
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
