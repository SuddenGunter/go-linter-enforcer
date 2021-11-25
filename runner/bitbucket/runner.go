package bitbucket

import (
	"context"
	"net/http"
	"time"

	"github.com/SuddenGunter/go-linter-enforcer/enforcer"
	"github.com/SuddenGunter/go-linter-enforcer/git"
	"github.com/SuddenGunter/go-linter-enforcer/repository"
	"go.uber.org/zap"
)

const (
	baseURL = "https://api.bitbucket.org"
	// pagelen allowed values range from 10 to 100.
	pagelen     = 50
	allowedLang = "go"
)

type Runner struct {
	gcp          *git.ClientProvider
	expectedFile []byte
	log          *zap.SugaredLogger
	cfg          Config
	client       *Client
}

func NewRunner(
	gcp *git.ClientProvider,
	expectedFile []byte,
	log *zap.SugaredLogger,
	cfg Config) *Runner {
	return &Runner{gcp: gcp, expectedFile: expectedFile, log: log, cfg: cfg, client: &Client{
		Client: http.Client{
			Timeout: 15 * time.Second,
		},
		Organization: cfg.Organization,
		Login:        cfg.Login,
		AppPassword:  cfg.AppPassword,
	}}
}

func (runner *Runner) Run(ctx context.Context) {
	timeout, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	repos, err := runner.client.LoadReposList(timeout)
	if err != nil {
		runner.log.Errorw("failed to get repositories list", "err", err, "organization", runner.cfg.Organization)
		return
	}

	for _, r := range repos {
		// todo: do concurrently
		enf := enforcer.NewEnforcer(runner.gcp, runner.log, repository.Author{
			Email: runner.cfg.Git.Email,
			Name:  runner.cfg.Git.Username,
		}, r, runner.expectedFile, runner.cfg.DryRun)

		branchName, err := enf.EnforceRules()
		if err != nil {
			runner.log.With("repo", r.Name).With("err", err).Error("failed to enforce rules, skipping repository")
			continue
		}

		runner.client.CreatePR(timeout, r, branchName)
	}
}
