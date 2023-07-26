package utils

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "127.0.0.1"
	database = "rental_system"
	port     = 5432
	user     = "postgres"
	password = "61926114"
)

func Connect_postgres() (DB *sql.DB, err error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, database,
	)
	db, err := sql.Open("postgres", connStr)
	if err == nil {
		return db, nil
	}
	return nil, err
}

func Execute_query(query string) (results *sql.Rows, err error) {
	conn, err := Connect_postgres()
	if err == nil {
		results, err := conn.Query(query)
		if err != nil {
			return nil, err
		} else {
			return results, nil
		}
	}
	conn.Close()
	return nil, err
}
