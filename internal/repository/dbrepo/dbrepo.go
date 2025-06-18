package dbrepo

import (
	"database/sql"

	"github.com/notirawap/doctor-registration/internal/repository"
)

// Create repo that has connection pool for Postgres
type postgresDBRepo struct {
	DB *sql.DB
}

// Initialize repo that includes Postgres connection pool
func NewPostgresRepo(conn *sql.DB) repository.DatabaseRepo {
	return &postgresDBRepo{
		DB: conn,
	}
}

// ------------------------------------------------------------------------------------------------------------------

// Create test database repository that does not connect to database at all swapped with the actual database repository
type testDBRepo struct {
	DB *sql.DB
}

func NewTestingRepo() repository.DatabaseRepo {
	return &testDBRepo{}
}
