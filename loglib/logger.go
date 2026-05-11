package loglib

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/lmittmann/tint"
)

var (
	ErrUnknownLevel = errors.New("unknown log level")
)

type SlogConfig struct {
	Level        string `yaml:"level" validate:"required"`
	OnlyStdout   bool   `yaml:"only_stdout"`
	LogFile      string `yaml:"log_file" validate:"required"`
	EnableCaller bool   `yaml:"enable_caller"`
}

type Logger interface {
	Warn(msg string, keysAndValues ...any)
	Info(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
	Debug(msg string, keysAndValues ...any)
	With(keysAndValues ...any) Logger
	WithGroup(name string) Logger
}

type Slog struct {
	onlyStdout bool
	closer     io.Closer
	closeOnce  *sync.Once
	slog       *slog.Logger
}

func NewSlog(cfg SlogConfig) (*Slog, error) {
	lvl, ok := stringToLevel(cfg.Level)
	if !ok {
		return nil, ErrUnknownLevel
	}
	consoleHandler := consoleHandler(cfg.EnableCaller, lvl)
	if cfg.OnlyStdout {
		s := slog.New(consoleHandler)
		return &Slog{
			onlyStdout: true,
			closer:     nil,
			closeOnce:  &sync.Once{},
			slog:       s,
		}, nil
	}
	fileHandler, closer, err := fileHandler(cfg.LogFile, cfg.EnableCaller, lvl)
	if err != nil {
		return nil, err
	}
	multiHandler := slog.NewMultiHandler(fileHandler, consoleHandler)
	s := slog.New(multiHandler)
	return &Slog{
		onlyStdout: false,
		closer:     closer,
		closeOnce:  &sync.Once{},
		slog:       s,
	}, nil
}

func (l *Slog) Info(msg string, args ...any) {
	l.slog.Info(msg, args...)
}
func (l *Slog) Error(msg string, args ...any) {
	l.slog.Error(msg, args...)
}
func (l *Slog) Debug(msg string, args ...any) {
	l.slog.Debug(msg, args...)
}
func (l *Slog) Warn(msg string, args ...any) {
	l.slog.Warn(msg, args...)
}

func (l *Slog) With(args ...any) Logger {
	return &Slog{
		slog: l.slog.With(args...),
	}
}
func (l *Slog) WithGroup(name string) Logger {
	return &Slog{
		slog: l.slog.WithGroup(name),
	}
}

func (l *Slog) Close() {
	if l.closer != nil {
		l.closeOnce.Do(func() {
			if err := l.closer.Close(); err != nil {
				panic(err)
			}
		})
		return
	}
	return
}

func consoleHandler(enableCaller bool, level slog.Level) slog.Handler {
	return tint.NewHandler(os.Stdout, &tint.Options{
		AddSource:  enableCaller,
		Level:      level,
		TimeFormat: time.DateTime,
	})
}

func fileHandler(fileName string, enableCaller bool, level slog.Level) (slog.Handler, io.Closer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return nil, nil, err
	}
	return slog.NewJSONHandler(file, &slog.HandlerOptions{
		AddSource: enableCaller,
		Level:     level,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && a.Value.Kind() == slog.KindTime {
				return slog.String(slog.TimeKey, a.Value.Time().Format(time.DateTime))
			}
			return a
		},
	}), file, nil
}

func stringToLevel(s string) (slog.Level, bool) {
	levelMap := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}
	level, ok := levelMap[strings.ToLower(s)]
	return level, ok
}
