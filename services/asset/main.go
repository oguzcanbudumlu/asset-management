package main

import (
	"asset-management/pkg/logger"
	"github.com/rs/zerolog"
)

func main() {
	logger.InitLogger(zerolog.InfoLevel)
	//log.Info().Msg("Hello from asset service")
}
