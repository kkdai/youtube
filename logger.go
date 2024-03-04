package youtube

import (
	"fmt"
	"log/slog"
	"os"
)

// The global logger for all Client instances
var Logger = getLogger(os.Getenv("LOGLEVEL"))

func SetLogLevel(value string) {
	Logger = getLogger(value)
}

func getLogger(logLevel string) *slog.Logger {
	levelVar := slog.LevelVar{}

	if logLevel != "" {
		if err := levelVar.UnmarshalText([]byte(logLevel)); err != nil {
			panic(fmt.Sprintf("Invalid log level %s: %v", logLevel, err))
		}
	}

	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: levelVar.Level(),
	}))
}
