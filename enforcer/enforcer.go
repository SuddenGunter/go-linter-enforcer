package enforcer

import (
	"github.com/SuddenGunter/go-linter-enforcer/repository"
	"go.uber.org/zap"
)

const (
	linterFileName = ".golangci.yaml"
	commitMessage  = "🤖 update " + linterFileName + " according to latest changes"
)

type GitClientProvider interface {
	OpenRepository(repo repository.Repository) (repository.GitClient, error)
}

type Enforcer struct {
	provider     GitClientProvider
	commitAuthor repository.Author
	log          *zap.SugaredLogger
	expectedFile []byte
	repo         repository.Repository
	dryRun       bool
}

func NewEnforcer(
	provider GitClientProvider,
	log *zap.SugaredLogger,

	commitAuthor repository.Author,
	repo repository.Repository,
	expectedFile []byte,

	dryRun bool) *Enforcer {
	return &Enforcer{
		provider:     provider,
		commitAuthor: commitAuthor,
		log:          log.With("repo", repo.Name),
		expectedFile: expectedFile,
		repo:         repo,
		dryRun:       dryRun,
	}
}

func (e *Enforcer) EnforceRules() {
	repo, err := e.provider.OpenRepository(e.repo)
	if err != nil {
		e.log.Errorw("errors when opening repository", "err", err)
		return
	}

	e.log.Debugw("repo opened")

	exists, err := repo.FileEquals(linterFileName, e.expectedFile)
	if err != nil {
		e.log.Debugw("errors when comparing file with existing", "err", err)
		return
	}

	if exists {
		e.log.Debugw("file exist and matches expected")
		return
	}

	e.log.Debugw("file doesn't match expected (or doesn't exist)")

	if err := repo.Replace(linterFileName, e.expectedFile); err != nil {
		e.log.Debugw("error when replacing file", "err", err)
		return
	}

	e.log.Debugw("replacing file")

	if e.dryRun {
		e.log.Debugw("dryRun mode enabled, no commits would be made")
		return
	}

	if err := repo.SaveChanges(commitMessage, e.commitAuthor); err != nil {
		e.log.Errorw("errors when trying to commit", "err", err)
		return
	}

	e.log.Debugw("committed new file", "file", linterFileName)
}
