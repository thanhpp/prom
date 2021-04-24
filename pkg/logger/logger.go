package logger

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/thanhpp/prom/pkg/logger/zaplog"
)

// EiLogger exported: DO NOT CALL THIS
type EiLogger iLogger

// iLogger ...
type iLogger interface {
	Fatal(message string)
	Fatalf(template string, args ...interface{})

	Panic(message string)
	Panicf(template string, args ...interface{})

	DPanic(message string)
	DPanicf(template string, args ...interface{})

	Error(message string)
	Errorf(template string, args ...interface{})

	Warn(message string)
	Warnf(template string, args ...interface{})

	Info(message string)
	Infof(template string, args ...interface{})

	Debug(message string)
	Debugf(template string, args ...interface{})
}

// Error
var (
	ErrLogDriverNotFound = errors.New("Dvergr. Log driver not found")
)

// loggerOption options for future use
type loggerOption interface {
}

// --------------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- FUNC GET SET ----------------------------------------------------------

var (
	iplmLogger iLogger
	_          iLogger = (*zaplog.ZapLogger)(nil) // zaplogger compile-time check
)

// NOTE: Maybe add write to file or send to log store later, options for future use
// Set ...
func Set(driver string, name string, enviroment string, level string, color bool, opts ...loggerOption) (err error) {
	switch strings.ToLower(driver) {
	case "zap":
		if err = zaplog.InitZapLogger(name, enviroment, level, color); err != nil {
			return err
		}
		iplmLogger = zaplog.GetZapLogger()

	default:
		return ErrLogDriverNotFound

	}
	return nil
}

// Get Must call Set(...)
func Get() (lg iLogger) {
	if iplmLogger == nil {
		_ = Set("zap", "default logger", "DEVELOPMENT", "DEBUG", true)
		iplmLogger.Info("Init default logger")
	}
	return iplmLogger
}

// -------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- UTILS ----------------------------------------------------------

// JsonIndentFmt format data to json indent for prettier print
func JsonIndentFmt(data interface{}) (jsIndent string, err error) {
	var (
		jsByte []byte
	)
	jsByte, err = json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}

	return string(jsByte), nil
}
