package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/notirawap/doctor-registration/internal/handlers"
)

func routes() http.Handler {
	// Create new router
	mux := chi.NewRouter()

	// Middleware 
	mux.Use(middleware.Recoverer) 
	mux.Use(NoSurf)               
	mux.Use(SessionLoad)         

	// Register a function to a URL path
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/doctor/form", handlers.Repo.DoctorForm)
	mux.Post("/doctor/form", handlers.Repo.PostDoctorForm)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/contact", handlers.Repo.Contact)

	mux.Get("/user/login", handlers.Repo.Login)
	mux.Post("/user/login", handlers.Repo.PostLogin)
	mux.Get("/user/logout", handlers.Repo.Logout)

	// Create protected routes only for login user
	mux.Route("/admin", func(mux chi.Router) {
		// apply middleware
		mux.Use(Auth) 
		// register subhandler to /admin/*
		mux.Get("/doctors-table", handlers.Repo.AdminDoctorsTable)
		mux.Get("/doctor-info/{id}", handlers.Repo.AdminShowDoctor)
		mux.Post("/update-doctor-info/{id}", handlers.Repo.AdminUpdateDoctor)
		mux.Get("/delete-doctor-info/{id}", handlers.Repo.AdminDeleteDoctor)
	})

	// Create means to serve the static files from the local file system into webserver-compatible format
	fileServer := http.FileServer(http.Dir("./static/"))             
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
