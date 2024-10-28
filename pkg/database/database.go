package database

import (
	"fmt"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type Database struct {
	Conn *gorm.DB
}

func NewDatabase(dbHost, dbPort, dbUser, dbPassword, dbName string) (*Database, error) {

	dsn := "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable"
	dsn = fmt.Sprintf(dsn, dbHost, dbPort, dbUser, dbPassword, dbName)
	gormLogger := logger.New(
		&log.Logger, // Use zerolog as the logger
		logger.Config{
			LogLevel:                  logger.Info,
			SlowThreshold:             200 * time.Millisecond,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLogger})

	if err != nil {
		log.Error().Err(err).Msg("Could not connect to the database")
		return nil, err
	}

	log.Info().Msg("Connected to the PostgreSQL database!")
	return &Database{Conn: conn}, nil
}

func (db *Database) Close() {
	sqlDB, err := db.Conn.DB()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get database instance")
		return
	}

	if err := sqlDB.Close(); err != nil {
		log.Error().Err(err).Msg("Failed to close the database connection")
	} else {
		log.Info().Msg("Database connection closed successfully")
	}
}
