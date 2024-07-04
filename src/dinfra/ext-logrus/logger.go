package extlogrus

import (
	"os"

	"github.com/AlgerDu/go-dream/src/dinfra"
	logrus "github.com/sirupsen/logrus"
)

type (
	LogLevel string

	LoggerOptions struct {
		MinLevel LogLevel
		File     *RotateFileOptions
	}
)

const (
	Trace LogLevel = "trace"
	Debug LogLevel = "debug"
	Info  LogLevel = "info"
	Warn  LogLevel = "warn"
	Fatal LogLevel = "fatal"
)

func levelToLogrus(level LogLevel) logrus.Level {
	switch level {
	case Trace:
		return logrus.TraceLevel
	case Debug:
		return logrus.DebugLevel
	case Info:
		return logrus.InfoLevel
	case Warn:
		return logrus.WarnLevel
	case Fatal:
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}

func NewDefaultOptions() *LoggerOptions {
	return &LoggerOptions{
		MinLevel: Info,
	}
}

func New(
	options *LoggerOptions,
) dinfra.Logger {

	if options == nil {
		options = NewDefaultOptions()
	}

	loger := logrus.New()

	minLevel := levelToLogrus(options.MinLevel)
	formatter := &logrus.TextFormatter{
		DisableColors:   true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	}

	loger.SetOutput(os.Stdout)
	loger.SetLevel(minLevel)
	loger.SetFormatter(formatter)

	if options.File != nil {

		fileOptions := options.File
		fileOptions.level = minLevel
		fileOptions.formatter = formatter

		fileHook, err := newRotateFileHook(fileOptions)
		if err != nil {
			panic(err)
		}

		loger.AddHook(fileHook)
	}

	return loger
}
