package repository

type Repository struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	MainBranch string `json:"mainBranch"`
}

type GitClient interface {
	FileEquals(path string, content []byte) (bool, error)
	Replace(path string, content []byte) error
	SaveChanges(commitMsg string, author Author) error
}
