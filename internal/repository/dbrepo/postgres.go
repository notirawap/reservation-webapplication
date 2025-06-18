package dbrepo // Define methods for postgresDBRepo type to fullfill repository.DatabaseRepo interface

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/notirawap/doctor-registration/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// Authenticate compares password the user typed against a hash of the password in the database
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Get user id and password from email
	var id int
	var hashedPassword string
	row := m.DB.QueryRowContext(ctx, "SELECT id, password FROM users WHERE email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	// If the email exists, bcrypt compares the hashed password from the database to the hashed string of what user typed
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password") // password mismatched
	} else if err != nil {
		return 0, "", err // other errors
	}
	return id, hashedPassword, nil
}

// InsertDoctor inserts a doctor record into the database
func (m *postgresDBRepo) InsertDoctor(doc models.Doctor) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO doctors (
			first_name, last_name, license, hospital, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := m.DB.ExecContext(ctx, stmt,
		doc.FirstName,
		doc.LastName,
		doc.License,
		doc.Hospital,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}

// GetDoctorByID returns a doctor record in the database
func (m *postgresDBRepo) GetDoctorByID(id int) (models.Doctor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var doc models.Doctor

	query := `
		SELECT id, first_name, last_name, license, hospital, created_at, updated_at
		FROM doctors
		WHERE id = $1
	`
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&doc.ID,
		&doc.FirstName,
		&doc.LastName,
		&doc.License,
		&doc.Hospital,
		&doc.CreatedAt,
		&doc.UpdatedAt,
	)
	if err != nil {
		return doc, err
	}
	return doc, nil
}

// UpdateDoctor updates a doctor record
func (m *postgresDBRepo) UpdateDoctorByID(doc models.Doctor) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE doctors SET first_name = $1, last_name = $2, license = $3, hospital = $4, updated_at = $5 WHERE id = $6`
	_, err := m.DB.ExecContext(ctx, query,
		doc.FirstName,
		doc.LastName,
		doc.License,
		doc.Hospital,
		time.Now(),
		doc.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

// DeleteDoctorByID deletes a doctor record in the database
func (m *postgresDBRepo) DeleteDoctorByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM doctors WHERE id = $1`
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

// AllDoctors returns a slice of all doctor records
func (m *postgresDBRepo) AllDoctors() ([]models.Doctor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var doctors []models.Doctor
	query := ` 
		SELECT id, first_name, last_name, license, hospital, created_at, updated_at
		FROM doctors
		ORDER BY license asc`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return doctors, err
	}
	defer rows.Close()
	for rows.Next() {
		var doc models.Doctor
		err := rows.Scan(
			&doc.ID,
			&doc.FirstName,
			&doc.LastName,
			&doc.License,
			&doc.Hospital,
			&doc.CreatedAt,
			&doc.UpdatedAt,
		)
		if err != nil {
			return doctors, err
		}
		doctors = append(doctors, doc)
	}
	if err = rows.Err(); err != nil {
		return doctors, err
	}
	return doctors, nil
}

// IsDoctorExists checks if there is an existing doctor information
func (m *postgresDBRepo) IsDoctorExist(firstName, lastName, license string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int

	// Check if license already exists
	query := `SELECT COUNT(*) FROM doctors WHERE license = $1`
	err := m.DB.QueryRowContext(ctx, query, license).Scan(&count)
	if err != nil {
		return fmt.Errorf("internal server error: database connection failed")
	}
	if count > 0 {
		return fmt.Errorf("doctor with this license already exists")
	}

	// Check if first name and last name already exist
	query = `SELECT COUNT(*) FROM doctors WHERE first_name = $1 AND last_name = $2`
	err = m.DB.QueryRowContext(ctx, query, firstName, lastName).Scan(&count)
	if err != nil {
		return fmt.Errorf("internal server error: database connection failed")
	}
	if count > 0 {
		return fmt.Errorf("doctor with the same name already exists")
	}

	return nil
}

