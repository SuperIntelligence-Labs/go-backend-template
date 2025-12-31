package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/SuperIntelligence-Labs/go-backend-template/internal/config"
)

var Log zerolog.Logger

func Init(level string) {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(parseLogLevel(level))

	var output io.Writer

	if config.IsDev() {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	} else {
		output = os.Stdout
	}

	// Build logger - add Caller() only in dev for performance
	logContext := zerolog.New(output).With().Timestamp()
	if config.IsDev() {
		logContext = logContext.Caller()
	}
	Log = logContext.Logger()
}

func parseLogLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}

func Debug() *zerolog.Event { return Log.Debug() }
func Info() *zerolog.Event  { return Log.Info() }
func Warn() *zerolog.Event  { return Log.Warn() }
func Error() *zerolog.Event { return Log.Error() }
func Fatal() *zerolog.Event { return Log.Fatal() }
