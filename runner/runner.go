package runner

import (
	"reflect"

	"go.uber.org/zap"
)

type Builder interface {
	CreateRunner(log *zap.SugaredLogger, config interface{}) Runner
	ConfigType() reflect.Type
}

type Runner interface {
	Run()
}
