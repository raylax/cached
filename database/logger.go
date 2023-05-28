package database

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
)

func newLoggerAdapter(logger hclog.Logger) Logger {
	return &loggerAdapter{
		logger: logger,
	}
}

type loggerAdapter struct {
	logger hclog.Logger
}

func (l *loggerAdapter) Errorf(s string, i ...interface{}) {
	if l.logger.IsError() {
		l.logger.Error(fmt.Sprintf(s, i...))
	}
}

func (l *loggerAdapter) Warningf(s string, i ...interface{}) {
	if l.logger.IsWarn() {
		l.logger.Warn(fmt.Sprintf(s, i...))
	}
}

func (l *loggerAdapter) Infof(s string, i ...interface{}) {
	if l.logger.IsInfo() {
		l.logger.Info(fmt.Sprintf(s, i...))
	}
}

func (l *loggerAdapter) Debugf(s string, i ...interface{}) {
	if l.logger.IsDebug() {
		l.logger.Debug(fmt.Sprintf(s, i...))
	}
}

type Logger interface {
	Errorf(string, ...interface{})
	Warningf(string, ...interface{})
	Infof(string, ...interface{})
	Debugf(string, ...interface{})
}
