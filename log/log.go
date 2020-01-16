package log

import (
	"errors"
	"fmt"
	"os"

	"github.com/gofier/framework/config"
	"github.com/gofier/framework/sentry"

	"github.com/sirupsen/logrus"
)

var l *logrus.Logger
var lLevel Level

func init() {
	l = logrus.New()
	l.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	l.Out = os.Stdout
}

func Init() {
	levelStr := config.GetString("app.log_level")
	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		panic(err)
	}
	lLevel = level
	l.SetLevel(lLevel)
}

func Error(msg ...interface{}) {
	fields := make(map[string]interface{})
	fields["level"] = "ERROR"
	sentry.CaptureError(errors.New(fmt.Sprintf("%v - %v", msg, fields)))
	l.Errorln(msg...)
}

func ErrorWithFields(msg interface{}, fields map[string]interface{}) {
	sentry.CaptureError(errors.New(fmt.Sprintf("%v - %v", msg, fields)))
	l.WithFields(fields).Errorln(msg)
}

func Warn(msg ...interface{}) {
	fields := make(map[string]interface{})
	fields["level"] = "WARN"
	sentry.CaptureMsg(fmt.Sprintf("%v", msg), fields)
	l.Warning(msg...)
}

func WarnWithFields(msg interface{}, fields map[string]interface{}) {
	sentry.CaptureMsg(fmt.Sprintf("%v", msg), fields)
	l.WithFields(fields).Warnln(msg)
}

func Info(msg ...interface{}) {
	l.Infoln(msg...)
}

func InfoWithFields(msg interface{}, fields map[string]interface{}) {
	l.WithFields(fields).Infoln(msg)
}

func Fatal(msg ...interface{}) {
	fields := make(map[string]interface{})
	fields["level"] = "FATAL"
	sentry.CaptureError(errors.New(fmt.Sprintf("%v - %v", msg, fields)))
	l.Fatalln(msg...)
}

func FatalWithFields(msg interface{}, fields map[string]interface{}) {
	sentry.CaptureError(errors.New(fmt.Sprintf("%v - %v", msg, fields)))
	l.WithFields(fields).Fatalln(msg)
}

func Debug(msg ...interface{}) {
	l.Debugln(msg...)
}

func DebugWithFields(msg interface{}, fields map[string]interface{}) {
	l.WithFields(fields).Debugln(msg)
}

func Trace(msg ...interface{}) {
	l.Traceln(msg...)
}

func TraceWithFields(msg interface{}, fields map[string]interface{}) {
	l.WithFields(fields).Traceln(msg)
}

func Panic(msg ...interface{}) {
	fields := make(map[string]interface{})
	fields["level"] = "PANIC"
	sentry.CaptureError(errors.New(fmt.Sprintf("%v - %v", msg, fields)))
	l.Panicln(msg...)
}

func PanicWithFields(msg interface{}, fields map[string]interface{}) {
	sentry.CaptureError(errors.New(fmt.Sprintf("%v - %v", msg, fields)))
	l.WithFields(fields).Panicln(msg)
}
