package repository

import (
	"github.com/notirawap/doctor-registration/internal/models"
)

type DatabaseRepo interface {
	Authenticate(email, testPassword string) (int, string, error)
	InsertDoctor(doc models.Doctor) error
	GetDoctorByID(id int) (models.Doctor, error)
	UpdateDoctorByID(doc models.Doctor) error
	DeleteDoctorByID(id int) error
	AllDoctors() ([]models.Doctor, error)
	IsDoctorExist(firstName, lastName, license string) error
}
