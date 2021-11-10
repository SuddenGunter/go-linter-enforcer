package repository

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage/memory"
)

func PushDemoBranch(auth transport.AuthMethod, repo Repository) error {
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:  repo.URL,
		Auth: auth,
	})
	if err != nil {
		return err
	}

	head, err := r.Head()
	if err != nil {
		return err
	}

	ref := plumbing.NewHashReference("refs/heads/demo-branch-1", head.Hash())

	err = r.Storer.SetReference(ref)
	if err != nil {
		return err
	}

	return r.Push(&git.PushOptions{
		Auth: auth,
	})
}
