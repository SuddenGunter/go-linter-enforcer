package config

import (
	"github.com/cristalhq/aconfig"
	"go.uber.org/zap"
)

type GitConfig struct {
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
	Email    string `env:"EMAIL"`
}

type Config struct {
	Git                  GitConfig `env:"GIT"`
	RepositoriesFile     string    `default:"repos.json" env:"REPOSITORIES_LIST_FILE"`
	ExpectedLinterConfig string    `default:"example.golangci.yaml" env:"LINTER_CONFIG_FILE"`
}

// FromEnv loads config values from env. Shuts down the application if something goes wrong.
func FromEnv(log *zap.SugaredLogger) Config {
	var cfg Config

	loader := aconfig.LoaderFor(&cfg, aconfig.Config{})
	if err := loader.Load(); err != nil {
		log.Fatalw("failed to parse config from env", "err", err)
	}

	return cfg
}
