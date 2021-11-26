package config

import (
	"github.com/cristalhq/aconfig"
	"go.uber.org/zap"
)

type GitConfig struct {
	Username              string `env:"USERNAME"`
	SSHPrivateKeyPath     string `env:"SSH_PRIVATE_KEY_PATH"`
	SSHPrivateKeyPassword string `env:"SSH_PRIVATE_KEY_PASSWORD"`
	Email                 string `env:"EMAIL"`
}

type Config struct {
	Git                  GitConfig `env:"GIT"`
	ExpectedLinterConfig string    `default:"example.golangci.yaml" env:"LINTER_CONFIG_FILE"`
	// todo: DryRun               bool      `default:"false" env:"DRY_RUN" flag:"dryRun"`
}

// todo
//func (cfg *Config) GetDryRunValue() bool {
//	return cfg.DryRun
//}
//
//type DryRunnable interface {
//	GetDryRunValue() bool
//}

// FromEnv loads config values from env. Shuts down the application if something goes wrong.
func FromEnv(log *zap.SugaredLogger, target interface{}) {
	loader := aconfig.LoaderFor(target, aconfig.Config{})
	if err := loader.Load(); err != nil {
		log.Fatalw("failed to parse config from env", "err", err)
	}

	// todo log.Debugw("running in mode", "dryRun", target.(DryRunnable).GetDryRunValue())
}
