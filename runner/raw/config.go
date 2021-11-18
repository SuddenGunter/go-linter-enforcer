package raw

import "github.com/SuddenGunter/go-linter-enforcer/config"

type Config struct {
	config.Config
	RepositoriesFile string `default:"repos.json" env:"REPOSITORIES_LIST_FILE"`
}
