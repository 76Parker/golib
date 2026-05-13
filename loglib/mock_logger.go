package loglib

import "log/slog"

type MockLogger struct {
	slog *slog.Logger
}

func (l *MockLogger) Info(msg string, args ...any) {
	l.slog.Info(msg, args...)
}
func (l *MockLogger) Error(msg string, args ...any) {
	l.slog.Error(msg, args...)
}
func (l *MockLogger) Debug(msg string, args ...any) {
	l.slog.Debug(msg, args...)
}
func (l *MockLogger) Warn(msg string, args ...any) {
	l.slog.Warn(msg, args...)
}

func (l *MockLogger) With(args ...any) Logger {
	return &MockLogger{
		slog: l.slog.With(args...),
	}
}
func (l *MockLogger) WithGroup(name string) Logger {
	return &Slog{
		slog: l.slog.WithGroup(name),
	}
}
