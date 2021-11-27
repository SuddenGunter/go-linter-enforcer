package bitbucket

import (
	"context"

	"go.uber.org/zap"

	"github.com/SuddenGunter/go-linter-enforcer/repository"
)

type DryRunAPIClient struct {
	realAPIClient APIClient
	log           *zap.SugaredLogger
}

func UseDryRun(realAPIClient APIClient, log *zap.SugaredLogger) APIClient {
	return &DryRunAPIClient{
		realAPIClient: realAPIClient,
		log:           log,
	}
}

func (c *DryRunAPIClient) LoadReposList(ctx context.Context) ([]repository.Repository, error) {
	return c.realAPIClient.LoadReposList(ctx)
}

func (c *DryRunAPIClient) CreatePR(ctx context.Context, r repository.Repository, name string) (CreatePRResponse, error) {
	c.log.With("repository", r.Name).With("branch", name).Debugw("tried to create PR for branch with dry run enabled. No PR would be created.")
	return CreatePRResponse{
		Links: struct {
			HTML linkWrapper `json:"html"`
		}{HTML: linkWrapper{
			Name: "NOP_LINK",
			Href: "NOP_LINK",
		}},
	}, nil
}
