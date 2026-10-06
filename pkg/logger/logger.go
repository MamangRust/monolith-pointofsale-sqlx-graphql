package logger

import (
	"os"
	"sync"

	gwlogger "github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/sdk/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerInterface is the logger contract shared across the point-of-sale services
// and pkg/resilience. It is a type alias of the payment-gateway logger interface
// so that resilience (which depends on that module) and this package expose the
// *same* type — passing a point-of-sale logger to a resilience guard compiles
// without any conversion.
type LoggerInterface = gwlogger.LoggerInterface

type Logger struct {
	Log *zap.Logger
}

var (
	once     sync.Once
	instance LoggerInterface
)

func NewLogger(service string, loggerProvider *log.LoggerProvider) (LoggerInterface, error) {
	var setupErr error

	once.Do(func() {
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		stdoutCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			zapcore.DebugLevel,
		)

		cores := []zapcore.Core{stdoutCore}
		if loggerProvider != nil {
			otelCore := otelzap.NewCore(
				service,
				otelzap.WithLoggerProvider(loggerProvider),
			)
			cores = append(cores, otelCore)
		}

		core := zapcore.NewTee(cores...)

		logger := zap.New(core,
			zap.AddCaller(),
			zap.AddCallerSkip(1),
			zap.AddStacktrace(zapcore.FatalLevel),
		)

		instance = &Logger{Log: logger}
	})

	return instance, setupErr
}

func (l *Logger) Info(message string, fields ...zap.Field) {
	l.Log.Info(message, fields...)
}

func (l *Logger) Fatal(message string, fields ...zap.Field) {
	l.Log.Fatal(message, fields...)
}

func (l *Logger) Debug(message string, fields ...zap.Field) {
	l.Log.Debug(message, fields...)
}

func (l *Logger) Error(message string, fields ...zap.Field) {
	l.Log.Error(message, fields...)
}

func (l *Logger) Warn(message string, fields ...zap.Field) {
	l.Log.Warn(message, fields...)
}

func (l *Logger) Check(level zapcore.Level, message string) *zapcore.CheckedEntry {
	return l.Log.Check(level, message)
}

func (l *Logger) With(fields ...zap.Field) LoggerInterface {
	return &Logger{Log: l.Log.With(fields...)}
}

func (l *Logger) Sync() error {
	return l.Log.Sync()
}

func GetInstance() LoggerInterface {
	return instance
}

func ResetInstance() {
	once = sync.Once{}
	instance = nil
}
