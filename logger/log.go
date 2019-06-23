// Package logger provides a convenient wrapper over zap.SugaredLogger
package logger

import (
	"go.uber.org/zap"
)

var logConf = zap.NewDevelopmentConfig()

// SetDebug sets the log level
func SetDebug(debug bool) {
	if debug {
		logConf.Level.SetLevel(zap.DebugLevel)
	} else {
		logConf.Level.SetLevel(zap.InfoLevel)
	}
}

// OutputToFile sets the log file location
func OutputToFile(file string) {
	logConf.OutputPaths = []string{
		file,
	}
}

// New creates a new zap.SugaredLogger from a global configuration
func New() (*zap.SugaredLogger, error) {
	log, logErr := logConf.Build()
	if logErr != nil {
		return nil, logErr
	}
	sugar := log.Sugar()
	return sugar, nil
}
