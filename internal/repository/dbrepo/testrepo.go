package dbrepo

// For testing, we don't want database function to access the database
// In unit test, we would like to run test without creating new database or running migration

import (
	"errors"

	"github.com/notirawap/doctor-registration/internal/models"
)

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	var id int
	var hashedPassword string
	return id, hashedPassword, nil
}

func (m *testDBRepo) InsertDoctor(doc models.Doctor) error {
	if doc.FirstName == "insert" {
		return errors.New("error")
	}
	return nil
}

func (m *testDBRepo) GetDoctorByID(id int) (models.Doctor, error) {
	var doc models.Doctor
	return doc, nil
}

func (m *testDBRepo) UpdateDoctorByID(doc models.Doctor) error {
	return nil
}

func (m *testDBRepo) DeleteDoctorByID(id int) error {
	return nil
}

func (m *testDBRepo) AllDoctors() ([]models.Doctor, error) {
	var doctors []models.Doctor
	return doctors, nil
}

func (m *testDBRepo) IsDoctorExist(firstName, lastName, license string) error {
	if firstName == "duplicate" && lastName == "duplicate" && license == "duplicate" {
		return errors.New("error")
	}
	return nil
}
