package enforcer

import (
	"errors"
	"fmt"

	"github.com/SuddenGunter/go-linter-enforcer/repository"
	"go.uber.org/zap"
)

const (
	LinterFileName = ".golangci.yml"
	CommitMessage  = "🤖 update " + LinterFileName + " according to latest changes"
)

var (
	ErrNothingToCommit = errors.New("nothing to commit: expected state matches actual")
	// todo: rename drymode -> dryrun
	ErrDryModeEnabled = errors.New("nothing to commit: dryMode enabled")
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

func (e *Enforcer) EnforceRules() (string, error) {
	repo, err := e.provider.OpenRepository(e.repo)
	if err != nil {
		return "", fmt.Errorf("errors when opening repository: %w", err)
	}

	e.log.Debugw("repo opened")

	equals, err := repo.FileEquals(LinterFileName, e.expectedFile)
	if err != nil {
		return "", fmt.Errorf("errors when comparing file with existing: %w", err)
	}

	if equals {
		return "", ErrNothingToCommit
	}

	e.log.Debugw("file doesn't match expected (or doesn't exist)")

	if err := repo.Replace(LinterFileName, e.expectedFile); err != nil {
		return "", fmt.Errorf("error when replacing file: %w", err)
	}

	e.log.Debugw("replacing file")

	if e.dryRun {
		return "", ErrDryModeEnabled
	}

	if err := repo.SaveChanges(CommitMessage, e.commitAuthor); err != nil {
		return "", fmt.Errorf("error when commit changes: %w", err)
	}

	e.log.Debugw("committed new file", "file", LinterFileName)

	return repo.CurrentBranchName()
}
