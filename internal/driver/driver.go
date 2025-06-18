// Setting up and managing a PostgreSQL database connection
package driver 

import (
	"database/sql"
	"time"

	// Import underlying driver (pgx)
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Create DB that holds the database connection pool
// DB struct can hold different database driver by adding or changing its member if you want to change driver at some point in the future
type DB struct {
	SQL *sql.DB
}

// Create Database Connection Pool for Postgres
func ConnectSQL(dsn string) (*DB, error) {
	// Create a new database for the application
	d, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}

	// test my connection
	if err = d.Ping(); err != nil {
		panic(err)
	}

	// Set attriutes for database connection pool
	d.SetMaxOpenConns(10)
	d.SetMaxIdleConns(5)
	d.SetConnMaxLifetime(5 * time.Minute)

	// Set database connection pool value to global variable
	var dbConn = &DB{} // holding connection pool
	dbConn.SQL = d

	// Ping database
	err = d.Ping()
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}
