package main

import (
	"os"
	"strings"

	"go.uber.org/zap"

	"github.com/SuddenGunter/go-linter-enforcer/runner/bitbucket"
	"github.com/SuddenGunter/go-linter-enforcer/runner/raw"

	"github.com/SuddenGunter/go-linter-enforcer/runner"

	"github.com/SuddenGunter/go-linter-enforcer/config"
	"github.com/SuddenGunter/go-linter-enforcer/logger"
)

const (
	Bitbucket = "BITBUCKET"
	Raw       = "RAW"
)

func main() {
	log := logger.Create()

	mode := os.Getenv("MODE")

	switch strings.ToUpper(mode) {
	case Bitbucket:
		launchRunner(log, bitbucket.RunnerBuilder{})
	case Raw:
		launchRunner(log, raw.RunnerBuilder{})
	default:
		log.With("err", "unknown mode").With("mode", mode).Fatal("failed to start runner")
	}
}

func launchRunner(log *zap.SugaredLogger, r runner.Builder) {
	cfg := r.Config()
	config.FromEnv(log, cfg)

	// todo: pass context for graceful shutdown
	r.CreateRunner(log, cfg).Run()
}
