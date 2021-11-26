package bitbucket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/SuddenGunter/go-linter-enforcer/enforcer"
	"github.com/SuddenGunter/go-linter-enforcer/repository"
)

type Client struct {
	http.Client
	Organization string
	Login        string
	AppPassword  string
}

func (c *Client) LoadReposList(ctx context.Context) ([]repository.Repository, error) {
	result := make([]repository.Repository, 0, 100)

	// todo: dry run support

	canMakeRequests := true

	for page := 1; canMakeRequests; page++ {
		url := fmt.Sprintf("%s/2.0/repositories/%s?page=%v&pagelen=%v", baseURL, c.Organization, page, pagelen)

		var repos getRepositoriesResponse

		err := c.performRequest(ctx, url, http.MethodGet, nil, &repos)
		if err != nil {
			return nil, fmt.Errorf("failed bitbucket api call: %w", err)
		}

		if repos.Next == "" {
			canMakeRequests = false
		}

		for _, r := range repos.Values {
			if r.Language == allowedLang {
				result = append(result, repository.Repository{
					Name:       r.Name,
					HTTPSURL:   r.Links.Self.Href,
					SSHURL:     c.getSSHURL(r.Links.Clone),
					MainBranch: r.Mainbranch.Name,
				})
			}
		}
	}

	return result, nil
}

type createPRRequest struct {
	Title  string `json:"title"`
	Source source `json:"source"`
}

type source struct {
	Branch branch `json:"branch"`
}

type branch struct {
	Name string `json:"name"`
}

type CreatePRResponse struct {
	Links struct {
		HTML linkWrapper `json:"html"`
	} `json:"links"`
}

func (c *Client) CreatePR(
	ctx context.Context,
	repo repository.Repository,
	branchName string) (CreatePRResponse, error) {

	// todo: dry run support
	url := fmt.Sprintf("%s/2.0/repositories/%s/%s/pullrequests", baseURL, c.Organization, repo.Name)
	req := createPRRequest{
		Title:  enforcer.CommitMessage,
		Source: source{branch{Name: branchName}},
	}

	body, err := json.Marshal(&req)
	if err != nil {
		return CreatePRResponse{}, fmt.Errorf("failed to create PR: %w", err)
	}

	var expectedResponse CreatePRResponse

	err = c.performRequest(ctx, url, http.MethodPost, bytes.NewBuffer(body), &expectedResponse)

	return expectedResponse, err
}

func (c *Client) performRequest(
	ctx context.Context,
	url, method string,
	body io.Reader,
	parsedResponse interface{}) error {
	request, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return fmt.Errorf("cannot form request. %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(c.Login, c.AppPassword)

	response, err := c.Do(request)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	// todo: what if body == nil?
	defer func() {
		if closeErr := response.Body.Close(); err == nil {
			err = fmt.Errorf("closing response body: %w", closeErr)
		}
	}()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("unable to read body: %w", err)
	}

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		if len(responseBody) == 0 {
			return fmt.Errorf("request failed but no detailed error received. status code: %v", response.StatusCode)
		}

		var apiErr map[string]interface{}
		if err = json.Unmarshal(responseBody, &apiErr); err != nil {
			return fmt.Errorf("failed unmarshal error form json body: %w", err)
		}

		return fmt.Errorf("api error: %v", apiErr)
	}

	if err = json.Unmarshal(responseBody, parsedResponse); err != nil {
		return fmt.Errorf("failed unmarshal response from JSON body: %w", err)
	}

	return nil
}

func (c *Client) getSSHURL(links []linkWrapper) string {
	for _, v := range links {
		if v.Name == "ssh" {
			return v.Href
		}
	}

	return ""
}

type getRepositoriesResponse struct {
	Values []struct {
		Links struct {
			Clone []linkWrapper `json:"clone"`
			Self  linkWrapper   `json:"self"`
		} `json:"links"`
		Name       string `json:"name"`
		Language   string `json:"language"`
		Mainbranch struct {
			Name string `json:"name"`
		} `json:"mainbranch"`
	} `json:"values"`
	Page int    `json:"page"`
	Next string `json:"next"`
}

type linkWrapper struct {
	Name string `json:"name"`
	Href string `json:"href"`
}
