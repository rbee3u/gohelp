package conflog

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger `env:"-"`
	ErrFile        string `env:""`
	File           string `env:""`
	Format         string `env:""`
	Level          string `env:""`
	ReportCaller   string `env:""`
}

type Option func(*Logger)

func WithErrFile(errFile string) Option {
	return func(logger *Logger) {
		logger.ErrFile = errFile
	}
}

func WithFile(file string) Option {
	return func(logger *Logger) {
		logger.File = file
	}
}

func WithFormat(format string) Option {
	return func(logger *Logger) {
		logger.Format = format
	}
}

func WithLevel(level string) Option {
	return func(logger *Logger) {
		logger.Level = level
	}
}

func WithReportCaller(reportCaller string) Option {
	return func(logger *Logger) {
		logger.ReportCaller = reportCaller
	}
}

func New(opts ...Option) (*Logger, error) {
	logger := &Logger{}

	if err := logger.SetDefaults(); err != nil {
		return nil, fmt.Errorf("failed to set defaults: %w", err)
	}

	for _, opt := range opts {
		opt(logger)
	}

	if err := logger.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize: %w", err)
	}

	return logger, nil
}

func (l *Logger) SetDefaults() error {
	if len(l.Format) == 0 {
		l.Format = "text"
	}

	if len(l.Level) == 0 {
		l.Level = logrus.InfoLevel.String()
	}

	if len(l.ReportCaller) == 0 {
		l.ReportCaller = "disable"
	}

	return nil
}

func (l *Logger) Initialize() error {
	l.Logger = logrus.New()

	callerPrettyfier := func(f *runtime.Frame) (string, string) {
		return "", fmt.Sprintf("%s:%d", f.Function, f.Line)
	}

	var formatter logrus.Formatter
	if strings.ToLower(l.Format) == "json" {
		formatter = &logrus.JSONFormatter{CallerPrettyfier: callerPrettyfier}
	} else {
		formatter = &logrus.TextFormatter{CallerPrettyfier: callerPrettyfier}
	}

	l.SetFormatter(formatter)

	if l.ReportCaller == "enable" {
		l.SetReportCaller(true)
	}

	if len(l.Level) != 0 {
		level, err := logrus.ParseLevel(l.Level)
		if err != nil {
			return fmt.Errorf("failed to parse level: %w", err)
		}

		l.SetLevel(level)
	}

	if len(l.File) != 0 {
		writer, err := os.OpenFile(l.File, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o666)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}

		l.SetOutput(writer)
	}

	if len(l.ErrFile) != 0 {
		errWriter, err := os.OpenFile(l.ErrFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o666)
		if err != nil {
			return fmt.Errorf("failed to open err file: %w", err)
		}

		l.AddHook(lfshook.NewHook(lfshook.WriterMap{
			logrus.ErrorLevel: errWriter,
			logrus.FatalLevel: errWriter,
			logrus.PanicLevel: errWriter,
		}, formatter))
	}

	return nil
}
