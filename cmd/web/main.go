package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/notirawap/doctor-registration/internal/config"
	"github.com/notirawap/doctor-registration/internal/driver"
	"github.com/notirawap/doctor-registration/internal/handlers"
	"github.com/notirawap/doctor-registration/internal/models"
	"github.com/notirawap/doctor-registration/internal/render"
)

const portNumber = ":8080"
const username = ""
const password = ""

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	// Start a webserver listening for the request in Go
	fmt.Printf("Starting application on port %v\n", portNumber)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(), // Set up Router
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}

var app config.AppConfig
var session *scs.SessionManager

func run() (*driver.DB, error) {
	// Connect to database
	fmt.Println("Connecting to database...")
	db, err := driver.ConnectSQL(fmt.Sprintf("host=localhost port=5432 dbname=meticuly_registration user=%s password=%s", username, password))
	if err != nil {
		log.Fatal("cannot connect to database.")
	}
	log.Println("Connected to database")

	// ************** Setting up App Config **************
	// Register models to put into the session
	gob.Register(models.User{})
	gob.Register(models.Doctor{})
	// Initialize session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true                 
	session.Cookie.SameSite = http.SameSiteLaxMode
	// Populate Session field in app.Config
	app.Session = session                          
	
	// Create new template cache
	tc, err := render.CreateTemplateCache() 
	if err != nil {
		return nil, err
	}
	// Populate TemplateCache in app.Config
	app.TemplateCache = tc

	// Create and set new repository for render package to access AppConfig
	render.NewRenderer(&app)

	// Create and set new repository for handlers package to access AppConfig and database
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	return db, nil
}
