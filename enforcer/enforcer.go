package enforcer

import (
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"go.uber.org/zap"

	"github.com/go-git/go-billy/v5/memfs"

	"github.com/SuddenGunter/go-linter-enforcer/repository"
	"github.com/beinan/fastid"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage/memory"
)

const (
	branchNameTemplate = "lintenforcer/2006-01-02-%v"
	linterFileName     = ".golangci.yaml"
)

type Enforcer struct {
	gitAuth      transport.AuthMethod
	log          *zap.SugaredLogger
	expectedFile []byte
}

func NewEnforcer(gitAuth transport.AuthMethod, log *zap.SugaredLogger, expectedFile []byte) *Enforcer {
	return &Enforcer{gitAuth: gitAuth, log: log, expectedFile: expectedFile}
}

func (e *Enforcer) EnforceRules(r repository.Repository) {
	repoLog := e.log.With("repo", r.Name)
	// get repo
	repo, err := e.getRepo(r)
	if err != nil {
		repoLog.Errorw("errors when opening repository", "err", err)
	}

	if err := e.checkIfFileIsTheSame(repo); err == nil {
		repoLog.Debugw("existing file matches expected file", "err", err)
		return
	}

	repoLog.Errorw("errors when comparing file with existing", "err", err)

	if err := e.tryReplaceFile(repo); err != nil {
		repoLog.Debugw("existing file matches expected file", "err", err)
	}

	if err := e.tryCommitChanges(repo); err != nil {
		repoLog.Errorw("errors when trying to commit", "err", err)
		return
	}

	// push new branch
	if err := repo.Push(&git.PushOptions{
		Auth: e.gitAuth,
	}); err != nil {
		repoLog.Errorw("errors when trying to push changes", "err", err)
		return
	}
}

func (e *Enforcer) getRepo(repo repository.Repository) (*git.Repository, error) {
	r, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL:  repo.URL,
		Auth: e.gitAuth,
	})
	if err != nil {
		return nil, err
	}

	head, err := r.Head()
	if err != nil {
		return nil, err
	}

	// create new branch to work on
	ref := plumbing.NewHashReference(e.getNewRefName(), head.Hash())
	if err = r.Storer.SetReference(ref); err != nil {
		return nil, err
	}

	return r, nil
}

func (e *Enforcer) getNewRefName() plumbing.ReferenceName {
	branchFormatWithoutTime := fmt.Sprintf(branchNameTemplate, fastid.CommonConfig.GenInt64ID())
	branchName := time.Now().UTC().Format(branchFormatWithoutTime)

	return plumbing.NewBranchReferenceName(branchName)
}

func (e *Enforcer) checkIfFileIsTheSame(repo *git.Repository) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	file, err := worktree.Filesystem.Open(linterFileName)
	if err != nil {
		return err
	}

	defer file.Close()

	existingFile, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	if len(existingFile) != len(e.expectedFile) {
		return errors.New("existing file length doesn't match expected")
	}

	for i, b := range existingFile {
		if b != e.expectedFile[i] {
			return errors.New("existing file doesn't match expected")
		}
	}

	return nil
}

func (e *Enforcer) tryReplaceFile(repo *git.Repository) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	if err = worktree.Filesystem.Remove(linterFileName); err != nil {
		return err
	}

	file, err := worktree.Filesystem.Open(linterFileName)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(e.expectedFile)
	if err != nil {
		return err
	}

	return nil
}

func (e *Enforcer) tryCommitChanges(repo *git.Repository) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	status, err := worktree.Status()
	if err != nil {
		return err
	}

	if status.IsClean() {
		return errors.New("nothing to commit")
	}

	return nil
}
