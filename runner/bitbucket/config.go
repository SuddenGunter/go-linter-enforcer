package bitbucket

import "github.com/SuddenGunter/go-linter-enforcer/config"

type Config struct {
	config.Config
	Login        string `env:"BITBUCKET_LOGIN"`
	AppPassword  string `env:"BITBUCKET_APP_PASSWORD"`
	Organization string `env:"BITBUCKET_ORGANIZATION"`
}
