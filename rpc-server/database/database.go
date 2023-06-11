package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	// Define the MySQL connection parameters
	dbUser = "root"
	dbPass = "password"
	dbHost = "mysql"
	dbPort = 3306
	dbName = "chat"

	migrationPath = "file://database/migration"

	retryInterval = 5 * time.Second
	timeout       = 90 * time.Second
)

func InitDB() (*sql.DB, error) {
	// Create the MySQL connection string
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	// Connect to the database
	db, err := connectDB(dataSourceName)
	if err != nil {
		return nil, err
	}
	log.Info().Msg("db connected")

	// Run the database migrations
	m, err := migrate.New(migrationPath, fmt.Sprintf("mysql://%s", dataSourceName))
	if err != nil {
		return nil, err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	log.Info().Msg("db migration successful")

	return db, nil
}

func connectDB(dataSourceName string) (*sql.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("unable to establish database connection within timeout")
		default:
			db, err := sql.Open("mysql", dataSourceName)
			if err != nil {
				log.Warn().Err(err).Msg("unable to open connection, retrying connection")
				time.Sleep(retryInterval)
				continue
			}

			if err = db.Ping(); err != nil {
				log.Warn().Err(err).Msg("unable to ping database, retrying connection")
				time.Sleep(retryInterval)
				continue
			}
			return db, nil
		}
	}
}
