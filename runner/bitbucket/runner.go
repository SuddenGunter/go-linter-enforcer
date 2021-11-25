package bitbucket

import (
	"context"

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
}

func NewRunner(
	gcp *git.ClientProvider,
	expectedFile []byte,
	log *zap.SugaredLogger,
	cfg Config) *Runner {
	return &Runner{gcp: gcp, expectedFile: expectedFile, log: log, cfg: cfg}
}

func (runner *Runner) Run() {
	// todo: pass ctx from caller
	repos, err := runner.loadReposList(context.Background())
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

		branchName, err := enf.EnforceRules() // todo: create PR from new branch to main
		if err != nil {
			runner.log.With("repo", r.Name).With("err", err).Error("failed to enforce rules, skipping repository")
			continue
		}

		// todo: pass ctx from caller
		runner.createPR(context.Background(), r, branchName)
	}
}

//nolint:gocognit,gocyclo
func (runner *Runner) loadReposList(ctx context.Context) ([]repository.Repository, error) {
	panic("not impl")
}

func (runner *Runner) createPR(ctx context.Context, r repository.Repository, name string) {
	// todo:
}
