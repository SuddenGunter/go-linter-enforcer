package git

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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
	branchNamePrefix   = "lintenforcer"
	branchNameTemplate = branchNamePrefix + "/2006-01-02-%v"
)

type ClientProvider struct {
	log  *zap.SugaredLogger
	auth transport.AuthMethod
}

func NewClientProvider(log *zap.SugaredLogger, auth transport.AuthMethod) *ClientProvider {
	return &ClientProvider{log: log, auth: auth}
}

func (p *ClientProvider) OpenRepository(repo repository.Repository) (repository.GitClient, error) {
	r, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL:  repo.SSHURL,
		Auth: p.auth,
	})
	if err != nil {
		return nil, err
	}

	// todo: check if 'lintenforcer/*' already exists
	iter, err := r.Branches()
	if err != nil {
		return nil, err
	}
	if err = iter.ForEach(func(reference *plumbing.Reference) error {
		if strings.Contains(reference.Name().Short(), branchNamePrefix) {
			return errors.New(branchNamePrefix + "/ branch already exists")
		}

		return nil
	}); err != nil {
		return nil, err
	}

	worktree, err := r.Worktree()
	if err != nil {
		return nil, err
	}

	if err = worktree.Checkout(&git.CheckoutOptions{
		Branch: p.getNewRefName(),
		Create: true,
	}); err != nil {
		return nil, err
	}

	return &Client{
		auth:       p.auth,
		log:        p.log,
		Repository: r,
		changed:    make(map[string]struct{}),
	}, nil
}

func (p *ClientProvider) getNewRefName() plumbing.ReferenceName {
	branchFormatWithoutTime := fmt.Sprintf(branchNameTemplate, fastid.CommonConfig.GenInt64ID())
	branchName := time.Now().UTC().Format(branchFormatWithoutTime)

	return plumbing.NewBranchReferenceName(branchName)
}

type Client struct {
	log     *zap.SugaredLogger
	changed map[string]struct{}
	*git.Repository
	auth transport.AuthMethod
}

// FileEquals compares existing file in path with content.
// Returns true if files are the same, or false if they are different.
// If file doesn't exist in the repository false is always returned, ignoring provided content.
func (c *Client) FileEquals(path string, content []byte) (bool, error) {
	worktree, err := c.Worktree()
	if err != nil {
		return false, err
	}

	file, err := worktree.Filesystem.Open(path)

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

	if len(existingFile) != len(content) {
		return false, nil
	}

	for i, b := range existingFile {
		if b != content[i] {
			return false, nil
		}
	}

	return true, nil
}

// Replace file and remember what file was replaced in the internal state.
func (c *Client) Replace(path string, content []byte) error {
	worktree, err := c.Worktree()
	if err != nil {
		return err
	}

	file, err := worktree.Filesystem.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return err
	}

	c.changed[path] = struct{}{}

	return nil
}

// SaveChanges creates commit and pushes it to remote origin.
func (c *Client) SaveChanges(commitMsg string, author repository.Author) error {
	worktree, err := c.Worktree()
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

	for f := range c.changed {
		_, err = worktree.Add(f)
		if err != nil {
			return err
		}
	}

	_, err = worktree.Commit(commitMsg, &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  author.Name,
			Email: author.Email,
			When:  time.Now().UTC(),
		},
		Committer: nil,
		Parents:   nil,
		SignKey:   nil,
	})

	if err != nil {
		return err
	}

	// push new branch
	return c.Push(&git.PushOptions{
		Auth: c.auth,
	})
}

// CurrentBranchName is returned in short form, e.g.: "refs/feature/foo" -> "feature/foo".
func (c *Client) CurrentBranchName() (string, error) {
	ref, err := c.Head()
	if err != nil {
		return "", err
	}

	return ref.Name().Short(), nil
}
