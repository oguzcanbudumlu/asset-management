package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/rs/zerolog/log"
)

type DatabaseRaw struct {
	Conn *sql.DB
}

func NewDatabaseRaw(host, port, user, password, dbname string) (*DatabaseRaw, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	log.Info().Msg(psqlInfo)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("Connected to the database!")
	return &DatabaseRaw{Conn: db}, nil
}

func (d *DatabaseRaw) Close() error {
	log.Info().Msg("Closing database connection...")
	return d.Conn.Close()
}
