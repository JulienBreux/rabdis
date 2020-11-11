package logger

// Logger represents the interface of logger
type Logger interface {
	AddDefaultField(Field)

	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Panic(msg string, fields ...Field)

	Infof(msg string, data ...interface{})
	Warningf(msg string, data ...interface{})
	Errorf(msg string, data ...interface{})
}

// Field represents the interface of logger field
type Field interface {
	Key() string
	Val() interface{}
}
