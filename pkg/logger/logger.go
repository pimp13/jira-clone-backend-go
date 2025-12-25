package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

type Logger = *zerolog.Logger

func New(isProduction bool) Logger {
	var output io.Writer

	if isProduction {
		output = os.Stdout
	} else {
		output = zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: "15:04:05.000",
			NoColor:    false,
		}
	}

	logger := zerolog.New(output).
		Level(zerolog.DebugLevel). // یا InfoLevel یا از env بخون
		With().
		Timestamp().
		Str("loc", "your-app").
		Logger()

	return &logger
}
