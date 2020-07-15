package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log constants
const (
	envLogLevel  = "LOG_LEVEL"
	envLogOutput = "LOG_OUTPUT"
)

// log variable
var (
	log logger
)

type logger struct {
	log *zap.Logger
}

func init() {
	logConfig := zap.Config{
		OutputPaths: []string{"stdout"},
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:     "level",
			TimeKey:      "time",
			MessageKey:   "msg",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	var err error
	if log.log, err = logConfig.Build(); err != nil {
		panic(err)
	}
}

// GetLogger function
// func GetLogger() {
// 	return log
// }

// func getLevel() zapcore.Level {
// 	switch strings.ToLower(strings.TrimSpace(os.Getenv(envLogLevel))) {
// 	case "debug":
// 		return zap.DebugLevel
// 	case "info":
// 		return zap.InfoLevel
// 	case "error":
// 		return zap.ErrorLevel
// 	default:
// 		return zap.InfoLevel
// 	}
// }

// func getOutput() string {
// 	output := strings.TrimSpace(os.Getenv(envLogOutput))
// 	if output == "" {
// 		return "stdout"
// 	}
// 	return output
// }

// func (l logger) Printf(format string, v ...interface{}) {
// 	if len(v) == 0 {
// 		Info(format)
// 	} else {
// 		Info(fmt.Sprintf(format, v...))
// 	}
// }

// func (l logger) Print(v ...interface{}) {
// 	Info(fmt.Sprintf("%v", v))
// }

// Info logger function
func Info(msg string, tags ...zap.Field) {
	log.log.Info(msg, tags...)
	log.log.Sync()
}

// Error logger function
func Error(msg string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("Error", err))
	log.log.Error(msg, tags...)
	log.log.Sync()
}
