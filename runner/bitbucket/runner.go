package bitbucket

import (
	"context"
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

type APIClient interface {
	LoadReposList(ctx context.Context) ([]repository.Repository, error)
	CreatePR(ctx context.Context, r repository.Repository, name string) (CreatePRResponse, error)
}

type Runner struct {
	gcp    *git.ClientProvider
	log    *zap.SugaredLogger
	client APIClient

	expectedFile []byte
	cfg          *Config
}

func NewRunner(
	gcp *git.ClientProvider,
	expectedFile []byte,
	log *zap.SugaredLogger,
	client APIClient,
	cfg *Config) *Runner {
	return &Runner{gcp: gcp, expectedFile: expectedFile, log: log, client: client, cfg: cfg}
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
		// todo: pass dryRun config
		enf := enforcer.NewEnforcer(runner.gcp, runner.log, repository.Author{
			Email: runner.cfg.Git.Email,
			Name:  runner.cfg.Git.Username,
		}, r, runner.expectedFile)

		branchName, err := enf.EnforceRules()
		if err != nil {
			runner.log.With("repo", r.Name).With("err", err).Error("failed to enforce rules, skipping repository")
			continue
		}

		resp, err := runner.client.CreatePR(timeout, r, branchName)
		if err != nil {
			runner.log.With("repo", r.Name).With("err", err).Error("failed to create PR")
			continue
		}

		runner.log.With("repo", r.Name).Debugw("created PR", "link", resp.Links.HTML)
	}
}
