package observability

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"syscall"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	OtelFieldBody                 = "Body"
	OtelFieldInstrumentationScope = "InstrumentationScope"
	OtelFieldSeverityText         = "SeverityText"
	OtelFieldTimestamp            = "Timestamp"
	OtelFieldTraceId              = "TraceId" //nolint:revive,stylecheck // Can't change the name of exported fields
	OtelFieldTraceFlags           = "TraceFlags"
	OtelFieldSpanId               = "SpanId" //nolint:revive,stylecheck // Can't change the name of exported fields
)

type Logger interface { //nolint:interfacebloat // This is a logging library, so we need to expose a lot of methods
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	WithFields(fields ...zap.Field) Logger
	WithError(err error) Logger
	WithContext(ctx context.Context) Logger
	Named(name string) Logger
	Sync() error
}

// logger instead.
type logger struct {
	config    *Config
	zapLogger *zap.Logger
	zapLevel  zap.AtomicLevel
}

func NewFromEnv() Logger {
	cfg := NewConfig()
	return processConfig(cfg)
}

func NewFromConfig(cfg *Config) Logger {
	if cfg == nil {
		return NewFromEnv()
	}

	if cfg.LogLevel == zapcore.InvalidLevel {
		cfg.LogLevel = zapcore.InfoLevel // Default to Info level if not set
	}

	return processConfig(cfg)
}

func (l *logger) Log(level zapcore.Level, msg string, fields ...zap.Field) {
	l.zapLogger.Log(level, msg, fields...)
}

func (l *logger) Debug(msg string, fields ...zap.Field) {
	l.zapLogger.Debug(msg, fields...)
}

func (l *logger) Info(msg string, fields ...zap.Field) {
	l.zapLogger.Info(msg, fields...)
}

func (l *logger) Warn(msg string, fields ...zap.Field) {
	l.zapLogger.Warn(msg, fields...)
}

func (l *logger) Error(msg string, fields ...zap.Field) {
	l.zapLogger.Error(msg, fields...)
}

func (l *logger) Fatal(msg string, fields ...zap.Field) {
	l.zapLogger.Fatal(msg, fields...)
}

func (l *logger) Sync() error {
	return ignoreSyncErorr(l.zapLogger.Sync())
}

func (l *logger) Named(name string) Logger {
	return &logger{
		config:    l.config,
		zapLogger: l.zapLogger.Named(name),
	}
}

func (l *logger) WithFields(fields ...zap.Field) Logger {
	return &logger{
		config:    l.config,
		zapLogger: l.zapLogger.With(fields...),
	}
}

func (l *logger) Level() zapcore.Level {
	return l.config.LogLevel
}

func (l *logger) WithContext(ctx context.Context) Logger {
	flds := make([]zap.Field, 0)

	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		spanContext := span.SpanContext()

		if spanContext.IsValid() {
			spanID := zap.String(OtelFieldSpanId, spanContext.SpanID().String())
			traceID := zap.String(OtelFieldTraceId, spanContext.TraceID().String())
			traceFlags := zap.Int(OtelFieldTraceFlags, int(spanContext.TraceFlags())) // byte converted to Int so it can be stored
			flds = append(flds, []zap.Field{spanID, traceID, traceFlags}...)
		} else {
			l.Debug("Found spanContext but is not valid")
		}
	} else {
		l.Debug("No active span found in the provided context")
	}

	return l.WithFields(flds...)
}

func (l *logger) WithError(err error) Logger {
	if err == nil {
		return l
	}

	return l.WithFields(extractErrorFields(err)...)
}

func extractErrorFields(err error) []zap.Field {
	f := make([]zap.Field, 0, 3)
	f = append(f, zap.String("exception.message", err.Error()), zap.String("exception.type", reflect.TypeOf(err).String()))

	// if the error carries a stack trace, add it to the fields
	switch errWithStack := err.(type) { //nolint:errorlint // We're not worried abut wrapped errors here, but do need the type info
	// match the interface of github.com/go-errors/errors
	case interface{ Stack() []byte }:
		f = append(f, zap.String("exception.stacktrace", string(errWithStack.Stack())))
	// match the interface of github.com/pkg/errors, as well as other 3rd-party
	// error libraries that use the convention of rendering stack traces via the Format method
	case fmt.Formatter:
		f = append(f, zap.String("exception.stacktrace", fmt.Sprintf("%+v", errWithStack)))
	}

	return f
}

func ignoreSyncErorr(err error) error {
	// Ignoring EINVAL for Linux and EBADF for Mac and ENOTTY for Unix
	if errors.Is(err, syscall.EINVAL) || errors.Is(err, syscall.EBADF) || errors.Is(err, syscall.ENOTTY) {
		return nil
	}
	return err
}

func processConfig(cfg *Config) Logger {
	zapLevel := zap.NewAtomicLevelAt(cfg.LogLevel)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.MessageKey = OtelFieldBody
	encoderConfig.TimeKey = OtelFieldTimestamp
	encoderConfig.LevelKey = OtelFieldSeverityText
	encoderConfig.NameKey = OtelFieldInstrumentationScope
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	enc := zapcore.NewJSONEncoder(encoderConfig)

	consolecore := zapcore.NewCore(enc, os.Stdout, zapLevel)

	zapLogger := zap.New(consolecore)

	return &logger{config: cfg, zapLogger: zapLogger, zapLevel: zapLevel}
}
