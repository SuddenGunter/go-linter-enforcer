package bitbucket

import (
	"github.com/SuddenGunter/go-linter-enforcer/enforcer"
	"github.com/SuddenGunter/go-linter-enforcer/git"
	"github.com/SuddenGunter/go-linter-enforcer/repository"
	"go.uber.org/zap"
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
	repos, err := runner.loadReposList()
	if err != nil {
		runner.log.Errorw("failed to get repositories list", "err", err, "organization", runner.cfg.Organization)
		return
	}

	for _, r := range repos {
		enf := enforcer.NewEnforcer(runner.gcp, runner.log, repository.Author{
			Email: runner.cfg.Git.Email,
			Name:  runner.cfg.Git.Username,
		}, r, runner.expectedFile, runner.cfg.DryRun)

		enf.EnforceRules()
	}
}

func (runner *Runner) loadReposList() ([]repository.Repository, error) {
	panic("implement me")
}
