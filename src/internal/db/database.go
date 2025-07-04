package db

import (
	"database/sql"
	_ "github.com/lib/pq"

	"fmt"
)

func Connection() (*sql.DB, error) {
	fmt.Println("Connecting to database ...")
	db, err := sql.Open("postgres", "postgres://admin:123@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("could not connect to postgres database: %w", err)
	}
	return db, nil
}
