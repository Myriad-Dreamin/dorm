package dorm

import (
	_ "github.com/go-sql-driver/mysql"
)

type Logger interface {
	Debug(msg string, keyvals ...interface{})
	Info(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})

	With(keyvals ...interface{}) Logger
}

type nopLogger struct{}

// NewNopLogger returns a logger that doesn't do anything.
func NewNopLogger() Logger { return &nopLogger{} }

func (*nopLogger) Info(string, ...interface{})  {}
func (*nopLogger) Debug(string, ...interface{}) {}
func (*nopLogger) Error(string, ...interface{}) {}

func (l *nopLogger) With(...interface{}) Logger {
	return l
}

var nl = new(nopLogger)
