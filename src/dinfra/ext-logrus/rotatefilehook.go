package extlogrus

import (
	"github.com/dounetio/slog/rotatefile"
	logrus "github.com/sirupsen/logrus"
)

type (
	RotateFileOptions struct {
		Filepath   string
		MaxSize    *uint64
		BackupNum  *uint
		BackupTime *uint

		level     logrus.Level
		formatter logrus.Formatter
	}

	rotateFileHook struct {
		options *RotateFileOptions
		writer  rotatefile.RotateWriter
	}
)

func newRotateFileHook(options *RotateFileOptions) (*rotateFileHook, error) {

	rfConfig, err := rotateFileOptionsToRotatefileConfig(options)
	if err != nil {
		return nil, err
	}

	writter, err := rfConfig.Create()
	if err != nil {
		return nil, err
	}

	return &rotateFileHook{
		options: options,
		writer:  writter,
	}, nil
}

func (hook *rotateFileHook) Levels() []logrus.Level {
	return logrus.AllLevels[:hook.options.level+1]
}

func (hook *rotateFileHook) Fire(entry *logrus.Entry) (err error) {
	b, err := hook.options.formatter.Format(entry)
	if err != nil {
		return err
	}
	hook.writer.Write(b)
	hook.writer.Flush()
	return nil
}

func rotateFileOptionsToRotatefileConfig(options *RotateFileOptions) (*rotatefile.Config, error) {

	config := rotatefile.NewConfig(options.Filepath)
	config.MaxSize = 2 * rotatefile.OneMByte

	if options.MaxSize != nil {
		config.MaxSize = *options.MaxSize
	}
	if options.BackupNum != nil {
		config.BackupNum = *options.BackupNum
	}
	if options.BackupTime != nil {
		config.BackupTime = *options.BackupTime
	}

	return config, nil
}
