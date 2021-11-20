package bitbucket

type getRepositoriesResponse struct {
	Values []struct {
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
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
