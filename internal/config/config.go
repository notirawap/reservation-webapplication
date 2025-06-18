package config

import (
	"html/template"

	"github.com/alexedwards/scs/v2"
)

// AppConfig holds the application configuration to be accessed from any part of the application.
type AppConfig struct {
	Session           *scs.SessionManager
	TemplateCache map[string]*template.Template
}
