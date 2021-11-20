package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/SuddenGunter/go-linter-enforcer/enforcer"
	"github.com/SuddenGunter/go-linter-enforcer/git"
	"github.com/SuddenGunter/go-linter-enforcer/repository"
	"go.uber.org/zap"
)

const (
	baseURL = "https://api.bitbucket.org"
	// pagelen allowed values range from 10 to 100.
	pagelen     = 50
	allowedLang = "go"
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
	repos, err := runner.loadReposList(context.Background())
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

func (runner *Runner) loadReposList(ctx context.Context) ([]repository.Repository, error) {
	httpClient := http.Client{
		Timeout: 15 * time.Second,
	}

	result := make([]repository.Repository, 0, 100)
	canMakeRequests := true
	page := 1
	for canMakeRequests {
		// todo: extract to bitbucketApiClient struct
		url := fmt.Sprintf("%s/2.0/repositories/%s?page=%v&pagelen=%v", baseURL, runner.cfg.Organization, page, pagelen)
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("cannot form request. %s", err)
		}

		request.Header.Set("Content-Type", "application/json")
		request.SetBasicAuth(runner.cfg.Login, runner.cfg.AppPassword)

		response, err := httpClient.Do(request)
		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}

		defer func() {
			if closeErr := response.Body.Close(); err == nil {
				err = fmt.Errorf("closing response body: %s", closeErr)
			}
		}()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("unable to read body: %s", err)
		}

		if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
			if len(body) == 0 {
				return nil, fmt.Errorf("request failed but no detailed error received. status code: %v", response.StatusCode)
			}

			var apiErr map[string]interface{}
			if err = json.Unmarshal(body, &apiErr); err != nil {
				return nil, fmt.Errorf("failed unmarshal error form json body: %w", err)
			}

			return nil, fmt.Errorf("api error: %w", apiErr)
		}

		var repos getRepositoriesResponse
		if err = json.Unmarshal(body, &repos); err != nil {
			return nil, fmt.Errorf("failed unmarshal response from JSON body: %w", err)
		}

		if repos.Next == "" {
			canMakeRequests = false
		}

		for _, r := range repos.Values {
			if r.Language == allowedLang {
				result = append(result, repository.Repository{
					Name:       r.Name,
					URL:        r.Links.Self.Href,
					MainBranch: r.Mainbranch.Name,
				})
			}
		}

		page++
	}

	return result, nil
}
