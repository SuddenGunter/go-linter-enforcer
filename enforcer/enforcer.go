package enforcer

import (
	"fmt"
	"time"

	"github.com/SuddenGunter/go-linter-enforcer/repository"
	"github.com/beinan/fastid"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage/memory"
)

const branchNameTemplate = "lintenforcer/2006-01-02-%v"

type Enforcer struct {
	GitAuth transport.AuthMethod
}

func (e *Enforcer) EnforceRules(repo repository.Repository) error {
	// get repo
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:  repo.URL,
		Auth: e.GitAuth,
	})
	if err != nil {
		return err
	}

	head, err := r.Head()
	if err != nil {
		return err
	}

	// create new branch
	ref := plumbing.NewHashReference(e.getNewRefName(), head.Hash())

	err = r.Storer.SetReference(ref)
	if err != nil {
		return err
	}

	// todo: compare existing linter config with expected HERE

	// push new branch
	return r.Push(&git.PushOptions{
		Auth: e.GitAuth,
	})
}

func (e *Enforcer) getNewRefName() plumbing.ReferenceName {
	branchFormatWithoutTime := fmt.Sprintf(branchNameTemplate, fastid.CommonConfig.GenInt64ID())
	branchName := time.Now().UTC().Format(branchFormatWithoutTime)

	return plumbing.NewBranchReferenceName(branchName)
}
