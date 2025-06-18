package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
	"github.com/notirawap/doctor-registration/internal/config"
	"github.com/notirawap/doctor-registration/internal/models"
	"github.com/notirawap/doctor-registration/internal/render"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"
var functions = template.FuncMap{}

// copy from main.go
func TestMain(m *testing.M) {
	gob.Register(models.User{})
	gob.Register(models.Doctor{})

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	app.Session = session

	tc, err := CreateTestTemplateCache()

	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc // Path is different

	repo := NewTestRepo(&app)
	NewHandlers(repo)

	render.NewRenderer(&app)

	os.Exit(m.Run())
}

// Copy from routes.go
func getRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	// mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/doctor/form", Repo.DoctorForm)
	mux.Post("/doctor/form", Repo.PostDoctorForm)
	mux.Get("/about", Repo.About)
	mux.Get("/contact", Repo.Contact)

	mux.Get("/user/login", Repo.Login)
	mux.Post("/user/login", Repo.PostLogin)
	mux.Get("/user/logout", Repo.Logout)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/doctors-table", Repo.AdminDoctorsTable)
		mux.Get("/doctor-info/{id}", Repo.AdminShowDoctor)
		mux.Post("/update-doctor-info/{id}", Repo.AdminUpdateDoctor)
		mux.Get("/delete-doctor-info/{id}", Repo.AdminDeleteDoctor)
	})
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}

// Copy from middleware.go
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.Session.Exists(r.Context(), "user_id") {
			session.Put(r.Context(), "error", "Please, Log in")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Copy from render package
func CreateTestTemplateCache() (map[string]*template.Template, error) {
	myCache := make(map[string]*template.Template)
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		templateSet, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		templateSet, err = templateSet.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}
		myCache[name] = templateSet
	}
	return myCache, nil
}
