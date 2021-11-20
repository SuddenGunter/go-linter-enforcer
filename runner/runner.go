package runner

import (
	"go.uber.org/zap"
)

type Builder interface {
	CreateRunner(log *zap.SugaredLogger, config interface{}) Runner
	Config() interface{}
}

type Runner interface {
	Run()
}
