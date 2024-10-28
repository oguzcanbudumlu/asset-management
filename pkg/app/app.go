package app

import (
	"asset-management/pkg/logger" // Logger import
	"context"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type App struct {
	fiber *fiber.App
}

func NewApp() *App {
	logger.InitLogger(zerolog.InfoLevel)
	return &App{
		fiber: fiber.New(),
	}
}

func (a *App) AddRoute(path string, handler fiber.Handler) {
	a.fiber.Get(path, handler)
}

func (a *App) Start(port string) {
	go func() {
		if err := a.fiber.Listen(port); err != nil {
			log.Fatal().Err(err).Msg("Error starting server")
		}
	}()

	a.setupGracefulShutdown(time.Second * 5)
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.fiber.ShutdownWithContext(ctx)
}

func (a *App) setupGracefulShutdown(timeout time.Duration) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := a.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msgf("Server forced to shutdown: %v", err)
	}

	log.Info().Msg("Server exited gracefully")
}
