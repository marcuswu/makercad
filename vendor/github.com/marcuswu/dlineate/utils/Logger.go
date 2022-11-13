package utils

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// var Logger = log.Logger.Level(zerolog.InfoLevel)
var Logger = log.Logger

func LogLevel() zerolog.Level {
	return Logger.GetLevel()
}
