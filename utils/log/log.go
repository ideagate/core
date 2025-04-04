package log

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func Debug(format string, args ...interface{}) {
	if len(args) == 0 {
		logrus.Debugln(format)
		return
	}
	logrus.Debugf(format, args...)
}

func Info(format string, args ...interface{}) {
	if len(args) == 0 {
		logrus.Infoln(format)
		return
	}
	logrus.Infof(format, args...)
}

func Warn(format string, args ...interface{}) {
	if len(args) == 0 {
		logrus.Warnln(format)
		return
	}
	logrus.Warnf(format, args...)
}

func Error(format string, args ...interface{}) {
	if len(args) == 0 {
		logrus.Errorln(format)
		return
	}
	logrus.Errorf(format, args...)
}

// ErrorP errors with prefix
func ErrorP(prefix string, args ...interface{}) {
	logrus.Error(fmt.Sprintf("[%s]:", prefix), args)
}

func Panic(format string, args ...interface{}) {
	if len(args) == 0 {
		logrus.Panicln(format)
	}
	logrus.Panicf(format, args...)
}

func Fatal(format string, args ...interface{}) {
	if len(args) == 0 {
		logrus.Fatalln(format)
	}
	logrus.Fatalf(format, args...)
}
