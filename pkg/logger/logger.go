package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func InitLogger(level zerolog.Level) {
	log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger().Level(level)
}
