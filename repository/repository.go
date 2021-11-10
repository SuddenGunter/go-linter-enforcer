package repository

type Repository struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	MainBranch string `json:"mainBranch"`
}
