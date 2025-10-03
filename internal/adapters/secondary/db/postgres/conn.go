package postgres

import (
	"database/sql"
	"os"

	"github.com/cheezecakee/logr"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostgresConn() *sql.DB {
	conn := os.Getenv("DB_CONN_STRING")
	if conn == "" {
		logr.Get().Errorf("DB_CONN_STRING environment variable is required")
		return nil
	}

	logr.Get().Infof("dbConn: %s", conn)

	db, err := sql.Open("pgx", conn)
	if err != nil {
		logr.Get().Errorf("Failed to connect to database: %v", err)
		return nil
	}

	if err := db.Ping(); err != nil {
		logr.Get().Errorf("Failed to ping database: %v", err)
		return nil

	}

	logr.Get().Info("Connected successfully!")
	return db
}
