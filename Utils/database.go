package utils

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
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

func Rows_iteration_error_check(rows *sql.Rows) error {
	if rows.Err() != nil {
		return rows.Err()
	}
	return nil
}

func Unique_constraint_violation_check(err error) bool {
	pqErr, _ := err.(*pq.Error)
	return pqErr.Code == "23505"
}

var ConnStr, ConnErr = Connect_postgres()
