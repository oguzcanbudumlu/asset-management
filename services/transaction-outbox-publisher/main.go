package main

import (
	"asset-management/internal/schedule"
	"asset-management/pkg/app"
	"asset-management/pkg/database"
	"asset-management/pkg/kafka"
	"asset-management/pkg/logger"
	"asset-management/services/transaction-outbox-publisher/publisher"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"os"
)

// @title Transaction Outbox Publisher API
// @version 1.0
// @description API documentation for Transaction Outbox Publisher Service.
// @BasePath /
func main() {
	logger.InitLogger(zerolog.InfoLevel)

	db, err := database.NewDatabaseRaw(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize database")
		return
	}
	appInstance := app.NewApp()
	appInstance.Fiber.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html")
	})

	appInstance.AddRoute("/swagger/*", fiberSwagger.WrapHandler)

	nextRepo := schedule.NewNextRepository(db.Conn)
	nextService := schedule.NewNextService(nextRepo)
	kafkaProducer := kafka.NewProducer(os.Getenv("KAFKA_BROKER"), os.Getenv("KAFKA_TOPIC"))

	s := publisher.NewService(nextService, kafkaProducer)
	c := publisher.NewController(s)

	appInstance.Fiber.Post("/trigger-publisher", c.TriggerPublisher)

	appInstance.Start(":8002")
	if dbCloseErr := db.Close(); dbCloseErr != nil {
		log.Error().Err(dbCloseErr).Msg("Failed to close database connection")
	}
}
