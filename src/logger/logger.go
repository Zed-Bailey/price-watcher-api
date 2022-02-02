package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

// set up zerolog with custom options
func SetupLogger() {
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	writer := io.Writer(logFile)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	Log = zerolog.New(writer).With().Timestamp().Logger()

}
