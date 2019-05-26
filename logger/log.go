package logger

import (
	"go.uber.org/zap"
)

var logConf = zap.NewDevelopmentConfig()

func SetDebug(debug bool) {
	if debug {
		logConf.Level.SetLevel(zap.DebugLevel)
	} else {
		logConf.Level.SetLevel(zap.InfoLevel)
	}
}

func New() (*zap.SugaredLogger, error) {
	log, logErr := logConf.Build()
	if logErr != nil {
		return nil, logErr
	}
	sugar := log.Sugar()
	return sugar, nil
}
