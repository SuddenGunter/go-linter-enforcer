package enforcer

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/SuddenGunter/go-linter-enforcer/repository"
	"github.com/beinan/fastid"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage/memory"
	"go.uber.org/zap"
)

const (
	branchNameTemplate = "lintenforcer/2006-01-02-%v"
	linterFileName     = ".golangci.yaml"
	commitMessage      = "🤖 update " + linterFileName + " according to latest changes"
)

type Author struct {
	Email string
	Name  string
}

type Enforcer struct {
	commitAuthor Author
	gitAuth      transport.AuthMethod
	log          *zap.SugaredLogger
	expectedFile []byte
}

func NewEnforcer(
	gitAuth transport.AuthMethod,
	commitAuthor Author,
	log *zap.SugaredLogger,
	expectedFile []byte) *Enforcer {
	return &Enforcer{gitAuth: gitAuth, commitAuthor: commitAuthor, log: log, expectedFile: expectedFile}
}

func (e *Enforcer) EnforceRules(r repository.Repository) {
	repoLog := e.log.With("repo", r.Name)
	// get repo
	repo, err := e.loadRepository(r)
	if err != nil {
		repoLog.Errorw("errors when opening repository", "err", err)
	}

	exists, err := e.checkIfFileIsTheSame(repo)
	if err != nil {
		repoLog.Debugw("errors when comparing file with existing", "err", err)
		return
	}

	if exists {
		repoLog.Debugw("file exist and matches expected")
		return
	}

	if err := e.tryReplaceFile(repo); err != nil {
		repoLog.Debugw("error when replacing file", "err", err)
		return
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

func (e *Enforcer) loadRepository(repo repository.Repository) (*git.Repository, error) {
	r, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL:  repo.URL,
		Auth: e.gitAuth,
	})
	if err != nil {
		return nil, err
	}

	worktree, err := r.Worktree()
	if err != nil {
		return nil, err
	}

	if err = worktree.Checkout(&git.CheckoutOptions{
		Branch: e.getNewRefName(),
		Create: true,
	}); err != nil {
		return nil, err
	}

	return r, nil
}

func (e *Enforcer) getNewRefName() plumbing.ReferenceName {
	branchFormatWithoutTime := fmt.Sprintf(branchNameTemplate, fastid.CommonConfig.GenInt64ID())
	branchName := time.Now().UTC().Format(branchFormatWithoutTime)

	return plumbing.NewBranchReferenceName(branchName)
}

func (e *Enforcer) checkIfFileIsTheSame(repo *git.Repository) (bool, error) {
	worktree, err := repo.Worktree()
	if err != nil {
		return false, err
	}

	file, err := worktree.Filesystem.Open(linterFileName)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, err
	}

	defer file.Close()

	existingFile, err := ioutil.ReadAll(file)
	if err != nil {
		return false, err
	}

	if len(existingFile) != len(e.expectedFile) {
		return false, nil
	}

	for i, b := range existingFile {
		if b != e.expectedFile[i] {
			return false, nil
		}
	}

	return true, nil
}

func (e *Enforcer) tryReplaceFile(repo *git.Repository) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	file, err := worktree.Filesystem.Create(linterFileName)
	if err != nil {
		return err
	}

	defer func() {
		file.Close()
	}()

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

	_, err = worktree.Add(linterFileName)
	if err != nil {
		return err
	}

	_, err = worktree.Commit(commitMessage, &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  e.commitAuthor.Name,
			Email: e.commitAuthor.Email,
			When:  time.Now().UTC(),
		},
		Committer: nil,
		Parents:   nil,
		SignKey:   nil,
	})

	return err
}
