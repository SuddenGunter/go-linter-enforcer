package raw

import (
	"github.com/SuddenGunter/go-linter-enforcer/enforcer"
	"github.com/SuddenGunter/go-linter-enforcer/git"
	"github.com/SuddenGunter/go-linter-enforcer/repository"
	"go.uber.org/zap"
)

type Runner struct {
	gcp          *git.ClientProvider
	repos        []repository.Repository
	expectedFile []byte
	log          *zap.SugaredLogger
	cfg          Config
}

func NewRunner(gcp *git.ClientProvider, repos []repository.Repository, expectedFile []byte, log *zap.SugaredLogger, cfg Config) *Runner {
	return &Runner{gcp: gcp, repos: repos, expectedFile: expectedFile, log: log, cfg: cfg}
}

func (runner *Runner) Run() {
	for _, r := range runner.repos {
		enf := enforcer.NewEnforcer(runner.gcp, runner.log, repository.Author{
			Email: runner.cfg.Git.Email,
			Name:  runner.cfg.Git.Username,
		}, r, runner.expectedFile, runner.cfg.DryRun)

		enf.EnforceRules()
	}
}
