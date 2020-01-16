package log

import "github.com/sirupsen/logrus"

type Level = logrus.Level

const (
	PANIC = logrus.PanicLevel
	FATAL = logrus.FatalLevel
	ERROR = logrus.ErrorLevel
	WARN  = logrus.WarnLevel
	INFO  = logrus.InfoLevel
	DEBUG = logrus.DebugLevel
	TRACE = logrus.TraceLevel
)
