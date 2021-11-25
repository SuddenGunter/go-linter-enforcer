package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
	canMakeRequests := true
	page := 1

	for canMakeRequests {
		url := fmt.Sprintf("%s/2.0/repositories/%s?page=%v&pagelen=%v", baseURL, c.Organization, page, pagelen)
		response, err := c.performRequest(ctx, url, http.MethodGet, nil)
		if err != nil {
			return nil, fmt.Errorf("failed api call: %w"), err)
		}

		request.Header.Set("Content-Type", "application/json")
		request.SetBasicAuth(c.Login, c.AppPassword)

		response, err := c.Do(request)
		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}

		defer func() {
			if closeErr := response.Body.Close(); err == nil {
				err = fmt.Errorf("closing response body: %w", closeErr)
			}
		}()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("unable to read body: %w", err)
		}

		if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
			if len(body) == 0 {
				return nil, fmt.Errorf("request failed but no detailed error received. status code: %v", response.StatusCode)
			}

			var apiErr map[string]interface{}
			if err = json.Unmarshal(body, &apiErr); err != nil {
				return nil, fmt.Errorf("failed unmarshal error form json body: %w", err)
			}

			return nil, fmt.Errorf("api error: %v", apiErr)
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
					HTTPSURL:   r.Links.Self.Href,
					SSHURL:     c.getSSHURL(r.Links.Clone),
					MainBranch: r.Mainbranch.Name,
				})
			}
		}

		page++
	}

	return result, nil
}

func (c *Client) performRequest(ctx context.Context, url, method string, body io.Reader) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("cannot form request. %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(c.Login, c.AppPassword)

	response, err := c.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	defer func() {
		if closeErr := response.Body.Close(); err == nil {
			err = fmt.Errorf("closing response body: %w", closeErr)
		}
	}()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read body: %w", err)
	}

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		if len(responseBody) == 0 {
			return nil, fmt.Errorf("request failed but no detailed error received. status code: %v", response.StatusCode)
		}

		var apiErr map[string]interface{}
		if err = json.Unmarshal(responseBody, &apiErr); err != nil {
			return nil, fmt.Errorf("failed unmarshal error form json body: %w", err)
		}

		return nil, fmt.Errorf("api error: %v", apiErr)
	}

	return responseBody, nil
}

func (c *Client) getSSHURL(links []LinkWrapper) string {
	for _, v := range links {
		if v.Name == "ssh" {
			return v.Href
		}
	}

	return ""
}
