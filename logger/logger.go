package logger

import (
	stdlog "log"

	"go.uber.org/zap"
)

func Create() *zap.SugaredLogger {
	log, err := zap.NewDevelopment()
	if err != nil {
		stdlog.Fatalln("failed to initialize logger")
	}

	return log.Sugar()
}
