package zaplog

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// --------------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- FUNC GET SET ----------------------------------------------------------

type ZapLogger struct{}

var (
	implZapLogger = new(ZapLogger)
)

// InitZapLogger create zap log for configuration
func InitZapLogger(name string, enviroment string, level string, color bool) (err error) {
	var (
		zapConfig zap.Config
	)

	// eviroment
	switch strings.ToUpper(enviroment) {
	case "PRODUCTION":
		zapConfig = zap.NewProductionConfig()
	default:
		zapConfig = zap.NewDevelopmentConfig()
	}

	// level
	switch strings.ToUpper(level) {
	case "FATAL":
		zapConfig.Level.SetLevel(zapcore.FatalLevel)
	case "PANIC":
		zapConfig.Level.SetLevel(zapcore.PanicLevel)
	case "DPANIC":
		zapConfig.Level.SetLevel(zapcore.DPanicLevel)
	case "ERROR":
		zapConfig.Level.SetLevel(zapcore.ErrorLevel)
	case "WARN":
		zapConfig.Level.SetLevel(zapcore.WarnLevel)
	case "INFO":
		zapConfig.Level.SetLevel(zapcore.InfoLevel)
	case "DEBUG":
		zapConfig.Level.SetLevel(zapcore.DebugLevel)
	default:
		zapConfig.Level.SetLevel(zapcore.InfoLevel)
	}

	// key
	zapConfig.EncoderConfig = zapcore.EncoderConfig{
		MessageKey: "message",

		NameKey:    "name",
		EncodeName: zapcore.FullNameEncoder,

		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,

		TimeKey:    "time",
		EncodeTime: zapcore.ISO8601TimeEncoder,

		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	// color
	if color {
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// build
	zlg, err := zapConfig.Build(zap.AddCallerSkip(1), zap.AddStacktrace(zap.DPanicLevel))
	if err != nil {
		return err
	}

	// replace global zap
	zap.ReplaceGlobals(zlg.Named(name))

	return nil
}

func GetZapLogger() (zlg *ZapLogger) {
	return implZapLogger
}

// ---------------------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- IMPLEMENT INTERFACE ----------------------------------------------------------

func (zlg ZapLogger) Fatal(message string) {
	zap.L().Fatal(message)
}

func (zlg ZapLogger) Fatalf(template string, args ...interface{}) {
	zap.S().Fatalf(template, args...)
}

func (zlg ZapLogger) Panic(message string) {
	zap.L().Panic(message)
}

func (zlg ZapLogger) Panicf(template string, args ...interface{}) {
	zap.S().Panicf(template, args...)
}

func (zlg ZapLogger) DPanic(message string) {
	zap.L().DPanic(message)
}

func (zlg ZapLogger) DPanicf(template string, args ...interface{}) {
	zap.S().DPanicf(template, args...)
}

func (zlg ZapLogger) Error(message string) {
	zap.L().Error(message)
}

func (zlg ZapLogger) Errorf(template string, args ...interface{}) {
	zap.S().Errorf(template, args...)
}

func (zlg ZapLogger) Warn(message string) {
	zap.L().Warn(message)
}

func (zlg ZapLogger) Warnf(template string, args ...interface{}) {
	zap.S().Warnf(template, args...)
}

func (zlg ZapLogger) Info(message string) {
	zap.L().Info(message)
}

func (zlg ZapLogger) Infof(template string, args ...interface{}) {
	zap.S().Infof(template, args...)
}

func (zlg ZapLogger) Debug(message string) {
	zap.L().Debug(message)
}

func (zlg ZapLogger) Debugf(template string, args ...interface{}) {
	zap.S().Debugf(template, args...)
}
