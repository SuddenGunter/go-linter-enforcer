package bitbucket

import (
	"reflect"

	"github.com/SuddenGunter/go-linter-enforcer/runner"
	"go.uber.org/zap"
)

type RunnerBuilder struct {
}

func (r RunnerBuilder) CreateRunner(log *zap.SugaredLogger, config interface{}) runner.Runner {
	panic("implement me")
}

func (r RunnerBuilder) ConfigType() reflect.Type {
	panic("implement me")
}

type Runner struct {
}

func (r *Runner) Run() {
	panic("implement me")
}
