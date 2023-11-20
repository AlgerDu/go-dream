package dinfra

import "github.com/sirupsen/logrus"

type (
	Logger interface {
		WithField(key string, value interface{}) *logrus.Entry
		WithFields(fields logrus.Fields) *logrus.Entry
		WithError(err error) *logrus.Entry

		Trace(args ...interface{})
		Debug(args ...interface{})
		Info(args ...interface{})
		Print(args ...interface{})
		Warn(args ...interface{})
		Warning(args ...interface{})
		Error(args ...interface{})
		Fatal(args ...interface{})
		Panic(args ...interface{})
	}
)

var (
	LogField_Source string = "Source"
	LogField_Track  string = "Track"
)
