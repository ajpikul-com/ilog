package ilog

import (
	"errors"
	"go.uber.org/zap"
	"os"
	"testing"
)

const (
	DEBUG int = iota
	INFO
	WARN
	ERROR
	NONE
)

const (
	RESET  = "\033[0m"
	RED    = "\033[31m"
	YELLOW = "\033[33m"
	GREEN  = "\033[32m"
)

// LogerInterface defines a simple interface to be used for logging. NOTE: Originally logr was used, but more features lead to less efficiency.
type LoggerInterface interface {
	Init() error
	Info(string)
	Error(string)
	Warn(string)
	Debug(string)
}

// EmptyLogger is a logger that can be used to turn off logging entirely.
type EmptyLogger struct{}

// Init is EmptyLogger's blank Init method
func (l *EmptyLogger) Init() error {
	return nil
}

// Info is EmptyLogger's blank Info method
func (l *EmptyLogger) Info(output string) {
}

// Error is EmptyLogger's blank Error method
func (l *EmptyLogger) Error(output string) {
}

// Warn is EmptyLogger's blank Error method
func (l *EmptyLogger) Warn(output string) {
}

// Debug is EmptyLogger's blank Error method
func (l *EmptyLogger) Debug(output string) {
}

// This statement ensures we're fufilling the interface we intend to
var _ LoggerInterface = &EmptyLogger{}

// SimpleLogger is a simple logger that writes to stderr or a path it's given. It is NOT safe for concurrent use.
type SimpleLogger struct {
	Path   string
	file   *os.File
	level  int
	prefix bool
}

func (l *SimpleLogger) Level(level int) {
	l.level = level
}

// Init attaches SimpleLogger to some path
func (l *SimpleLogger) Init() error {
	var err error
	if len(l.Path) == 0 {
		l.file = os.Stderr
		l.prefix = true
	} else {
		l.file, err = os.OpenFile(l.Path, os.O_APPEND|os.O_WRONLY, 0644)
		l.prefix = false
	}
	return err
}

// Info writes the info string to the output for SimpleLogger
func (l *SimpleLogger) Info(output string) {
	if l.level > INFO {
		return
	}
	_, _ = l.file.WriteString(output + "\n")
}

// Error writes the error string to the output for SimpleLogger
func (l *SimpleLogger) Error(output string) {
	if l.level > ERROR {
		return
	}
	if l.prefix {
		_, _ = l.file.WriteString(RED + output + RESET + "\n")
		return
	}
	_, _ = l.file.WriteString(output + "\n")
}

func (l *SimpleLogger) Warn(output string) {
	if l.level > WARN {
		return
	}
	if l.prefix {
		_, _ = l.file.WriteString(YELLOW + output + RESET + "\n")
		return
	}
	_, _ = l.file.WriteString(output + "\n")
}

func (l *SimpleLogger) Debug(output string) {
	if l.level > DEBUG {
		return
	}
	if l.prefix {
		_, _ = l.file.WriteString(GREEN + output + RESET + "\n")
		return
	}
	_, _ = l.file.WriteString(output + "\n")
}

var _ LoggerInterface = &SimpleLogger{}

// ZapWrap produces a uber-zap logging connection
type ZapWrap struct {
	// Sugar is a flag to indicate whether we should use a Sugared logger
	Sugar bool
	// Paths lets us set the logging paths, otherwise we use stderr
	Paths []string
	// Level is a debugging level
	Level int
	// ZapLogger is the underlying ZapLogger
	ZapLogger *zap.Logger
	// SugarLogger is the underlying SugaredLogger
	SugarLogger *zap.SugaredLogger
	// infoFunc is the function called by Info() method
	infoFunc func(output string)
	// errorFunc is the function called by the Error() method
	errorFunc func(output string)
	debugFunc func(output string)
	warnFunc  func(output string)
}

// Init starts a production level zap logger, which we use since we don't use all the same logging levels as Zap. It will switch the info or error func depending on whether or not its a sugared logger
func (z *ZapWrap) Init() error {
	config := zap.NewProductionConfig()
	if len(z.Paths) > 0 {
		config.OutputPaths = z.Paths
		switch z.Level {
		case DEBUG:
			config.Level.SetLevel(zap.DebugLevel)
		case INFO:
			config.Level.SetLevel(zap.InfoLevel)
		case WARN:
			config.Level.SetLevel(zap.WarnLevel)
		case ERROR:
			config.Level.SetLevel(zap.ErrorLevel)
		case NONE:
			config.Level.SetLevel(zap.DPanicLevel)
		}
	}
	z.ZapLogger, _ = config.Build(zap.AddCallerSkip(2)) // Can add more callers here with zap.AddCaller
	if z.Sugar {
		z.SugarLogger = z.ZapLogger.Sugar()
		z.infoFunc = func(output string) {
			z.SugarLogger.Info(output)
		}
		z.errorFunc = func(output string) {
			z.SugarLogger.Error(output)
		}
		z.debugFunc = func(output string) {
			z.SugarLogger.Debug(output)
		}
		z.warnFunc = func(output string) {
			z.SugarLogger.Warn(output)
		}
	} else {
		z.infoFunc = func(output string) {
			z.ZapLogger.Info(output)
		}
		z.errorFunc = func(output string) {
			z.ZapLogger.Error(output)
		}
		z.debugFunc = func(output string) {
			z.Debug(output)
		}
		z.warnFunc = func(output string) {
			z.Warn(output)
		}
	}
	return nil
}

// Info is ZapWraps Info method, but just a wrapper for z.infoFunc
func (z *ZapWrap) Info(output string) {
	z.infoFunc(output)
}

// Error is ZapWraps Error method, just a wrapper for z.errorFunc
func (z *ZapWrap) Error(output string) {
	z.errorFunc(output)
}

func (z *ZapWrap) Warn(output string) {
	z.warnFunc(output)
}

func (z *ZapWrap) Debug(output string) {
	z.debugFunc(output)
}

var _ LoggerInterface = &ZapWrap{}

// testLogger is an implmentation of a LoggerInterface but, like the "testing" package, requires one to set either the benchmark or testing variable.
type testLogger struct {
	t testing.TB
}

// NewTestLogger is a wrapper to create a logger for use with those tests
func NewTestLogger(t testing.TB) testLogger {
	ret := testLogger{t: t}
	ret.Init()
	return ret
}

// Info supplies the less severe testing method of testLogger
func (t *testLogger) Info(output string) {
	t.t.Log(output)
}

// Error supplies the most severe testing method of testLogger
func (t *testLogger) Error(output string) {
	t.t.Error(output)
}

func (t *testLogger) Warn(output string) {
	t.t.Error(output)
}

func (t *testLogger) Debug(output string) {
	t.t.Log(output)
}

// Init is the required Init function of a LoggerInterface for testLogger
func (t *testLogger) Init() error {
	if t.t == nil {
		return errors.New("You must set testLogger.t")
	}
	return nil
}

// This just ensures we're fufilling the interface we intend to at compile time
var _ LoggerInterface = &testLogger{}
