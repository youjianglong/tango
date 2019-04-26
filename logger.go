// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tango

import (
	"io"
	"time"

	"log"
)

const (
	LogLevelError = iota
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
)

// Logger defines the logger interface for tango use
type Logger interface {
	Debugf(format string, v ...interface{})
	Debug(v ...interface{})
	Debugln(v ...interface{})
	Infof(format string, v ...interface{})
	Info(v ...interface{})
	Infoln(v ...interface{})
	Warnf(format string, v ...interface{})
	Warn(v ...interface{})
	Warnln(v ...interface{})
	Errorf(format string, v ...interface{})
	Error(v ...interface{})
	Errorln(v ...interface{})
}

// CompositeLogger defines a composite loggers
type CompositeLogger struct {
	loggers []Logger
}

// NewCompositeLogger creates a composite loggers
func NewCompositeLogger(logs ...Logger) Logger {
	return &CompositeLogger{loggers: logs}
}

// Debugf implementes Logger interface
func (l *CompositeLogger) Debugf(format string, v ...interface{}) {
	for _, l := range l.loggers {
		l.Debugf(format, v...)
	}
}

// Debug implementes Logger interface
func (l *CompositeLogger) Debug(v ...interface{}) {
	for _, l := range l.loggers {
		l.Debug(v...)
	}
}

// Debug implementes Logger interface
func (l *CompositeLogger) Debugln(v ...interface{}) {
	for _, l := range l.loggers {
		l.Debugln(v...)
	}
}

// Infof implementes Logger interface
func (l *CompositeLogger) Infof(format string, v ...interface{}) {
	for _, l := range l.loggers {
		l.Infof(format, v...)
	}
}

// Info implementes Logger interface
func (l *CompositeLogger) Info(v ...interface{}) {
	for _, l := range l.loggers {
		l.Info(v...)
	}
}

// Infoln implementes Logger interface
func (l *CompositeLogger) Infoln(v ...interface{}) {
	for _, l := range l.loggers {
		l.Infoln(v...)
	}
}

// Warnf implementes Logger interface
func (l *CompositeLogger) Warnf(format string, v ...interface{}) {
	for _, l := range l.loggers {
		l.Warnf(format, v...)
	}
}

// Warn implementes Logger interface
func (l *CompositeLogger) Warn(v ...interface{}) {
	for _, l := range l.loggers {
		l.Warn(v...)
	}
}

// Warnln implementes Logger interface
func (l *CompositeLogger) Warnln(v ...interface{}) {
	for _, l := range l.loggers {
		l.Warnln(v...)
	}
}

// Errorf implementes Logger interface
func (l *CompositeLogger) Errorf(format string, v ...interface{}) {
	for _, l := range l.loggers {
		l.Errorf(format, v...)
	}
}

// Error implementes Logger interface
func (l *CompositeLogger) Error(v ...interface{}) {
	for _, l := range l.loggers {
		l.Error(v...)
	}
}

// Errorln implementes Logger interface
func (l *CompositeLogger) Errorln(v ...interface{}) {
	for _, l := range l.loggers {
		l.Errorln(v...)
	}
}

type DefaultLogger struct {
	*log.Logger
	level int
}

func (t *DefaultLogger) SetLevel(level int) {
	t.level = level
}

func (t *DefaultLogger) Debug(v ...interface{}) {
	if t.level >= LogLevelDebug {
		t.Logger.Println(append([]interface{}{" DEBUG "}, v...)...)
	}
}

func (t *DefaultLogger) Debugf(format string, v ...interface{}) {
	if t.level >= LogLevelDebug {
		t.Logger.Printf(" DEBUG "+format, v...)
	}
}

func (t *DefaultLogger) Debugln(v ...interface{}) {
	if t.level >= LogLevelDebug {
		t.Logger.Println(append([]interface{}{"DEBUG"}, v...)...)
	}
}

func (t *DefaultLogger) Info(v ...interface{}) {
	if t.level >= LogLevelInfo {
		t.Logger.Println(append([]interface{}{" INFO "}, v...)...)
	}
}

func (t *DefaultLogger) Infof(format string, v ...interface{}) {
	if t.level >= LogLevelInfo {
		t.Logger.Printf(" INFO "+format, v...)
	}
}

func (t *DefaultLogger) Infoln(v ...interface{}) {
	if t.level >= LogLevelInfo {
		t.Logger.Println(append([]interface{}{"INFO"}, v...)...)
	}
}

func (t *DefaultLogger) Warn(v ...interface{}) {
	if t.level >= LogLevelWarn {
		t.Logger.Println(append([]interface{}{" WARN "}, v...)...)
	}
}

func (t *DefaultLogger) Warnf(format string, v ...interface{}) {
	if t.level >= LogLevelWarn {
		t.Logger.Printf(" WARN "+format, v...)
	}
}

func (t *DefaultLogger) Warnln(v ...interface{}) {
	if t.level >= LogLevelWarn {
		t.Logger.Println(append([]interface{}{"WARN"}, v...)...)
	}
}

func (t *DefaultLogger) Error(v ...interface{}) {
	if t.level >= LogLevelError {
		t.Logger.Println(append([]interface{}{" ERROR "}, v...)...)
	}
}

func (t *DefaultLogger) Errorf(format string, v ...interface{}) {
	if t.level >= LogLevelError {
		t.Logger.Printf(" ERROR "+format, v...)
	}
}

func (t *DefaultLogger) Errorln(v ...interface{}) {
	if t.level >= LogLevelError {
		t.Logger.Println(append([]interface{}{"ERROR"}, v...)...)
	}
}

// NewLogger use the default logger with special writer
func NewDefaultLogger(out io.Writer) *DefaultLogger {
	l := new(DefaultLogger)
	l.Logger = log.New(out, "[tango] ", log.LstdFlags|log.Lshortfile|log.Ltime)
	l.level = LogLevelInfo
	return l
}

// LogInterface defines logger interface to inject logger to struct
type LogInterface interface {
	SetLogger(Logger)
}

// Log implementes LogInterface
type Log struct {
	Logger
}

// SetLogger implementes LogInterface
func (l *Log) SetLogger(log Logger) {
	l.Logger = log
}

// Logging returns handler to log informations
func Logging() HandlerFunc {
	return func(ctx *Context) {
		start := time.Now()
		p := ctx.Req().URL.Path
		if len(ctx.Req().URL.RawQuery) > 0 {
			p = p + "?" + ctx.Req().URL.RawQuery
		}

		ctx.Debugln("Started", ctx.Req().Method, p, "for", ctx.IP())

		if action := ctx.Action(); action != nil {
			if l, ok := action.(LogInterface); ok {
				l.SetLogger(ctx.Logger)
			}
		}

		ctx.Next()

		if !ctx.Written() {
			if ctx.Result == nil {
				ctx.Result = NotFound()
			}
			ctx.HandleError()
		}

		statusCode := ctx.Status()

		if statusCode >= 200 && statusCode < 400 {
			ctx.Infoln(ctx.Req().Method, statusCode, time.Since(start), p)
		} else {
			ctx.Errorln(ctx.Req().Method, statusCode, time.Since(start), p, ctx.Result)
		}
	}
}
