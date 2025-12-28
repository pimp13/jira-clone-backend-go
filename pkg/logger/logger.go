package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

type Logger = *zerolog.Logger

func New(isProduction bool) Logger {
	var output io.Writer
	var logger zerolog.Logger

	if isProduction {
		output = os.Stdout
	} else {
		output = zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: "15:04:05.000",
			NoColor:    false,
		}
	}

	logger = zerolog.New(output).
		Level(zerolog.DebugLevel). // یا InfoLevel یا از env بخون
		With().
		Timestamp().
		Caller().
		// Str("loc", "your-app").
		Logger()

	return &logger
}

// func New(loc string) Logger {
// 	var output io.Writer

// 	output = zerolog.ConsoleWriter{
// 		Out:        os.Stderr,
// 		TimeFormat: "15:04:05.000",
// 		NoColor:    false,
// 	}

// 	logger := zerolog.New(output).
// 		Level(zerolog.DebugLevel). // یا InfoLevel یا از env بخون
// 		With().
// 		Timestamp().
// 		Caller().
// 		Str("loc", loc).
// 		Logger()

// 	return &logger
// }
