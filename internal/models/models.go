// store all database models
package models

import (
	"time"
)

type User struct {
	ID        int
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Doctor struct {
	ID        int
	FirstName string
	LastName  string
	License   string
	Hospital  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
