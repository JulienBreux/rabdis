package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"

	stackdriver "github.com/TV4/logrus-stackdriver-formatter"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

const (
	formatterText        = "text"
	formatterJSON        = "json"
	formatterStackDriver = "stackdriver"
)

// logger represents the internal package logger
type logger struct {
	o *Options
	l *logrus.Logger

	dfs []Field
}

// New creates a new logger instance
func New(opts ...Option) (Logger, error) {
	o, err := newOptions(opts...)
	if err != nil {
		return nil, err
	}

	// Create logger
	l := logrus.New()

	// Set level
	if o.Level != "" {
		if err := l.Level.UnmarshalText([]byte(o.Level)); err != nil {
			return nil, err
		}
	}

	// Set formatter
	l.Formatter = convertFormatter(o)

	lgr := &logger{
		o: o,
		l: l,

		dfs: o.DefaultFields,
	}

	return lgr, nil
}

// AddDefaultFields adds default fields into to logger
func (l *logger) AddDefaultField(f Field) {
	l.dfs = append(l.dfs, f)
}

// Debug logs at debug level
func (l *logger) Debug(msg string, fields ...Field) {
	l.fields(fields...).Debug(msg)
}

// Info logs at info level
func (l *logger) Info(msg string, fields ...Field) {
	l.fields(fields...).Info(msg)
}

// Warn logs at warn level
func (l *logger) Warn(msg string, fields ...Field) {
	l.fields(fields...).Warn(msg)
}

// Error logs at error level
func (l *logger) Error(msg string, fields ...Field) {
	l.fields(fields...).Error(msg)
}

// Fatal logs at fatal level
func (l *logger) Fatal(msg string, fields ...Field) {
	l.fields(fields...).Fatal(msg)
}

// Panic logs at panic level
func (l *logger) Panic(msg string, fields ...Field) {
	l.fields(fields...).Panic(msg)
}

// Infof logs at info with formated message
func (l *logger) Infof(msg string, data ...interface{}) {
	l.Info(fmt.Sprintf(msg, data...))
}

// Warningf logs at warning with formated message
func (l *logger) Warningf(msg string, data ...interface{}) {
	l.Warn(fmt.Sprintf(msg, data...))
}

// Errorf logs at error with formated message
func (l *logger) Errorf(msg string, data ...interface{}) {
	l.Error(fmt.Sprintf(msg, data...))
}

func (l *logger) fields(fields ...Field) *logrus.Entry {
	if f := l.serializeField(fields...); f != nil {
		return l.l.WithFields(f)
	}

	return logrus.NewEntry(l.l)
}

// serializeField from logger field
func (l *logger) serializeField(fields ...Field) logrus.Fields {
	fs := logrus.Fields{}

	// Default fields
	if len(l.o.DefaultFields) != 0 && l.o.Format != formatterStackDriver {
		for _, f := range l.o.DefaultFields {
			fs[f.Key()] = f.Val()
		}
	}

	// Normal fields
	if len(fields) != 0 {
		for _, f := range fields {
			fs[f.Key()] = f.Val()
		}
	}

	return fs
}

// convertFormatter from string to logger Fromatter
func convertFormatter(opts *Options) logrus.Formatter {
	switch opts.Format {
	case formatterText:
		return new(prefixed.TextFormatter)
	case formatterJSON:
		return new(logrus.JSONFormatter)
	case formatterStackDriver:
		return stackdriver.NewFormatter(
			stackdriver.WithService(opts.InstName),
			stackdriver.WithVersion(opts.InstVersion),
		)
	default:
		return new(logrus.TextFormatter)
	}
}
