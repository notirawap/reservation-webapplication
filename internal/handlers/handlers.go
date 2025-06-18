package handlers

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/notirawap/doctor-registration/internal/config"
	"github.com/notirawap/doctor-registration/internal/driver"
	"github.com/notirawap/doctor-registration/internal/forms"
	"github.com/notirawap/doctor-registration/internal/models"
	"github.com/notirawap/doctor-registration/internal/render"
	"github.com/notirawap/doctor-registration/internal/repository"
	"github.com/notirawap/doctor-registration/internal/repository/dbrepo"
)

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo // handlers can connect to database via database connection pool repo (that is not specific to certain database driver)
}

var Repo *Repository // all the handlers can access to the repository

// 1. create a new repository populated with AppConfig and Database Connection Pool
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL),
	}
}

func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(), // does not return the database for testing
	}
}

// 2. set the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) DoctorForm(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "doctor-form.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostDoctorForm(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Internal Server Error")
		log.Printf("%s\n%s", r.URL.Path, debug.Stack())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Get the values from the HTML form element
	var newDoc models.Doctor

	newDoc.FirstName = r.Form.Get("first_name")
	newDoc.LastName = r.Form.Get("last_name")
	newDoc.License = r.Form.Get("license")
	newDoc.Hospital = r.Form.Get("hospital")

	// Server Side Form Validation
	form := forms.New(r.PostForm)

	// Populate Errors field by appending error message if error occurs
	form.Required("first_name", "last_name", "license", "hospital")
	form.NumLength("license", 5)

	// If error occurs, re-render the page and re-populate information that the user previously entered.
	if !form.Valid() {
		// Create stringmap storing dates in string as template data
		data := make(map[string]interface{})
		data["doctor"] = newDoc

		// re-render the template
		render.Template(w, r, "doctor-form.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// Call database function to check for duplicates
	err = m.DB.IsDoctorExist(newDoc.FirstName, newDoc.LastName, newDoc.License)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", fmt.Sprintf("%s", err))
		http.Redirect(w, r, "/doctor/form", http.StatusTemporaryRedirect)
		return
	}

	// If not duplicate, insert data into database
	err = m.DB.InsertDoctor(newDoc)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Internal Server Error")
		log.Printf("%s\n%s", r.URL.Path, debug.Stack())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "You have successfully completed the form")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	// Generates a new session token
	_ = m.App.Session.RenewToken(r.Context())

	// Parse Form
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Internal Server Error")
		log.Printf("%s\n%s", r.URL.Path, debug.Stack())
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// Form validation
	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	// Authenticate
	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		log.Printf("%s\n%s", r.URL.Path, debug.Stack())
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// After successfully authenticated, log in by storing the user id into the session
	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/admin/doctors-table", http.StatusSeeOther)
}

func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())
	m.App.Session.Put(r.Context(), "flash", "Logged out")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (m *Repository) AdminDoctorsTable(w http.ResponseWriter, r *http.Request) {
	doctors, err := m.DB.AllDoctors()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Internal Server Error")
		log.Printf("%s\n%s", r.URL.Path, debug.Stack())
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	data := make(map[string]interface{})
	data["doctors"] = doctors
	render.Template(w, r, "admin-doctors-table.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminShowDoctor(w http.ResponseWriter, r *http.Request) {
	// Extract URL parameters
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[3])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Internal Server Error")
		log.Printf("%s\n%s", r.URL.Path, debug.Stack())
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Get the doctor information from the database
	doctor, err := m.DB.GetDoctorByID(id)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Internal Server Error")
		log.Printf("%s\n%s", r.URL.Path, debug.Stack())
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Store doctor information into session
	m.App.Session.Put(r.Context(), "doctor", doctor)

	// Store doctor information into template data
	data := make(map[string]interface{})
	data["doctor"] = doctor
	render.Template(w, r, "admin-doctor-show.page.tmpl", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

func (m *Repository) AdminUpdateDoctor(w http.ResponseWriter, r *http.Request) {
	// Get value from the form
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Internal Server Error")
		log.Printf("%s\n%s", r.URL.Path, debug.Stack())
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	doctor := m.App.Session.Get(r.Context(), "doctor").(models.Doctor)
	doctor.FirstName = r.Form.Get("first_name")
	doctor.LastName = r.Form.Get("last_name")
	doctor.License = r.Form.Get("license")
	doctor.Hospital = r.Form.Get("hospital")

	// Server Side Form Validation
	form := forms.New(r.PostForm)

	// Populate Errors field by appending error message if error occurs
	form.Required("first_name", "last_name", "license", "hospital")
	form.NumLength("license", 5)

	// If error occurs, re-render the page and re-populate information that the user previously entered.
	if !form.Valid() {
		// re-render the template
		data := make(map[string]interface{})
		data["doctor"] = doctor
		render.Template(w, r, "admin-doctor-show.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// Call database function to check for duplicates
	err = m.DB.IsDoctorExist(doctor.FirstName, doctor.LastName, doctor.License)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", fmt.Sprintf("%s",err))
		http.Redirect(w, r, "/admin/doctors-table", http.StatusSeeOther)
		return
	}

	// Update the doctor model
	err = m.DB.UpdateDoctorByID(doctor)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Internal Server Error")
		log.Printf("%s\n%s", r.URL.Path, debug.Stack())
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	m.App.Session.Put(r.Context(), "flash", "Changes Saved")

	// remove doctor information from session
	m.App.Session.Remove(r.Context(), "doctor")

	// Redirect to the dashboard
	http.Redirect(w, r, "/admin/doctors-table", http.StatusSeeOther)
}

func (m *Repository) AdminDeleteDoctor(w http.ResponseWriter, r *http.Request) {
	// get doctor information from session
	doctor := m.App.Session.Get(r.Context(), "doctor").(models.Doctor)

	// Delete the doctor information from the database
	err := m.DB.DeleteDoctorByID(doctor.ID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Internal Server Error")
		log.Printf("%s\n%s", r.URL.Path, debug.Stack())
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	m.App.Session.Put(r.Context(), "flash", fmt.Sprintf("%s  | %s | license: %s | is successfully deleted", doctor.FirstName, doctor.LastName, doctor.License))

	// remove doctor information from session
	m.App.Session.Remove(r.Context(), "doctor")

	// Redirect to the dashboard
	http.Redirect(w, r, "/admin/doctors-table", http.StatusSeeOther)
}
