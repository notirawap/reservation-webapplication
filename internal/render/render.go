package render

// Render html template as a Golang template
import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/notirawap/doctor-registration/internal/config"
	"github.com/notirawap/doctor-registration/internal/models"
)

// template.FuncMap specifies functions to be accessed in Golang Template
var functions = template.FuncMap{}
var pathToTemplates = "./templates"
var app *config.AppConfig

// NewRenderer sets the config for the render package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// Determine the template data to be available on the page
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	// Retrieves the value from the session then delete it immediately (for one-time messages to the user)
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")

	// Passing default CSRF token as template data
	td.CSRFToken = nosurf.Token(r)

	// Check if the user is authenticated
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = true
	}

	return td
}

// Render html template from cache
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	// Get the requested template set from cache
	t, ok := app.TemplateCache[tmpl]
	if !ok {
		return fmt.Errorf("template %s does not exist", tmpl)
	}
	// Render the template to buffer
	buf := new(bytes.Buffer)
	td = AddDefaultData(td, r)
	err := t.Execute(buf, td) // Replace template variables with corresponding td fields
	if err != nil {
		return err
	}

	// Write the rendered template from the buffer to the HTTP response writer
	_, err = buf.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}

// Store parsed templates in cache (map)
func CreateTemplateCache() (map[string]*template.Template, error) {
	// Create map as an empty template cache
	myCache := make(map[string]*template.Template)

	// Get all template files named *.pages.tmpl in ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	// Range through all template files
	// When rendering a template using a layout, parse the template then its associated layouts.
	for _, page := range pages {
		// 1. Parses the page file into this new template set
		name := filepath.Base(page)
		templateSet, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		// 2. Loads and parses all layout files and adds their defined templates to the current template set
		templateSet, err = templateSet.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}
		// Store the parsed template set in the cache
		myCache[name] = templateSet
	}
	return myCache, nil
}
