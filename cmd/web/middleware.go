package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// NoSerf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	// Create a new CSRFHandler
	csrfHandler := nosurf.New(next)
	// Set the base cookie to use when building a CSRF token cookie
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",              // apply to the entire site
		Secure:   false, // true if running on https or in production
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// Auth is applied to the protected routes for logged in user
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		// If not authenticated, redirect to /user/login page. 
		if !app.Session.Exists(r.Context(), "user_id") {
			session.Put(r.Context(), "error", "Please, Log in")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		// If authenticated, access the protected routes
		next.ServeHTTP(w, r)
	})
}
