package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"

	"fmt"
)

func Connection() (*sql.DB, error) {
	fmt.Println("Connecting to database ...")
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://admin:123@localhost:5432/postgres?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("could not connect to postgres database: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("could not ping postgres database: %w", err)
	}
	return db, nil
}
