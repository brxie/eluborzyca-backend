package ilog

import (
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/brxie/eluborzyca-backend/config"
	log "github.com/sirupsen/logrus"
)

// verb is log level
var verb = map[string]log.Level{
	"PANIC": log.PanicLevel,
	"FATAL": log.FatalLevel,
	"WARN":  log.WarnLevel,
	"INFO":  log.InfoLevel,
	"DEBUG": log.DebugLevel,
	"TRACE": log.TraceLevel,
}
var logger *log.Logger
var once sync.Once

// singleton
func getLogger() *log.Logger {
	once.Do(func() {
		logger = log.New()
	})
	return logger
}

// Init set log configurations.
// Configures verbosity if passed in config. If not set, uses default.
func init() {
	formatter := new(log.TextFormatter)
	formatter.TimestampFormat = "02-01-2006 15:04:05"

	logger = getLogger()
	logger.SetFormatter(formatter)

	levelStr := strings.ToUpper(config.Viper.GetString("LOG_LEVEL"))
	if v, ok := verb[levelStr]; ok {
		logger.SetLevel(v)
	}
}

// Trace for log
func Trace(start time.Time, message ...interface{}) {
	if pc, file, line, ok := runtime.Caller(1); ok {
		getLogger().WithField("file", filepath.Base(file)+":"+strconv.Itoa(line)).
			WithField("func", runtime.FuncForPC(pc).Name()).
			WithField("duration", time.Since(start)).
			Trace(message...)
	} else {
		getLogger().Trace(message...)
	}
}

// Debug for log
func Debug(message ...interface{}) {
	if pc, file, line, ok := runtime.Caller(1); ok {
		getLogger().WithField("file", filepath.Base(file)+":"+strconv.Itoa(line)).WithField("func", runtime.FuncForPC(pc).Name()).Debug(message...)
	} else {
		getLogger().Debug(message...)
	}
}

// Info for log
func Info(message ...interface{}) {
	if pc, file, line, ok := runtime.Caller(1); ok {
		getLogger().WithField("file", filepath.Base(file)+":"+strconv.Itoa(line)).WithField("func", runtime.FuncForPC(pc).Name()).Info(message...)
	} else {
		getLogger().Info(message...)
	}
}

// Warn for log
func Warn(message ...interface{}) {
	if pc, file, line, ok := runtime.Caller(1); ok {
		getLogger().WithField("file", filepath.Base(file)+":"+strconv.Itoa(line)).WithField("func", runtime.FuncForPC(pc).Name()).Warn(message...)
	} else {
		getLogger().Warn(message...)
	}
}

// Error for log
func Error(message ...interface{}) {
	if pc, file, line, ok := runtime.Caller(1); ok {
		getLogger().WithField("file", filepath.Base(file)+":"+strconv.Itoa(line)).WithField("func", runtime.FuncForPC(pc).Name()).Error(message...)
	} else {
		getLogger().Error(message...)
	}
}

// Fatal for log
func Fatal(message ...interface{}) {
	if pc, file, line, ok := runtime.Caller(1); ok {
		getLogger().WithField("file", filepath.Base(file)+":"+strconv.Itoa(line)).WithField("func", runtime.FuncForPC(pc).Name()).Fatal(message...)
	} else {
		getLogger().Fatal(message...)
	}
}

// Panic for log
func Panic(message ...interface{}) {
	if pc, file, line, ok := runtime.Caller(1); ok {
		getLogger().WithField("file", filepath.Base(file)+":"+strconv.Itoa(line)).WithField("func", runtime.FuncForPC(pc).Name()).Panic(message...)
	} else {
		getLogger().Panic(message...)
	}
}
