package render

import (
	"net/http"
	"testing"

	"github.com/notirawap/doctor-registration/internal/models"
)

func getSession() (*http.Request, error) {
	// Create a new request
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}
	// Put session data into the context
	ctx := r.Context()                                    // Get the context
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session")) // updates the ctx to include session data; returns context containing the session state.
	// Update the request object r to include the new context added with session data
	r = r.WithContext(ctx)
	return r, nil
}

func TestNewTemplates(t *testing.T) {
	NewRenderer(app)
}

func TestAddDefaultData(t *testing.T) {
	// Create models.TemplateData (argument)
	var td models.TemplateData
	// Create a new request that has Session data in its context (argument)
	r, err := getSession()
	if err != nil {
		t.Error("failed")
	}
	// 1. test for valid response
	session.Put(r.Context(), "flash", "123") // Put session test data into the context
	result := AddDefaultData(&td, r)
	if result.Flash != "123" {
		t.Error("flash value of 123 not found in session")
	}

	// 2. test for invalid response
	session.Remove(r.Context(), "flash") // Remove session test data into the context
	result = AddDefaultData(&td, r)
	if result.Flash == "123" {
		t.Error("flash value of 123 found in session")
	}
}

func TestTemplate(t *testing.T) {
	// Create template cache (argument)
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
	app.TemplateCache = tc
	// Create http.Request (argument)
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	// Create http.ResponseWriter (argument)
	var ww myWriter // a type that satisfied the ResponseWriter interface

	// 1. test for valid response
	err = Template(&ww, r, "about.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("error writing template to browser")
	}
	// 2. test for invalid response
	err = Template(&ww, r, "non-existent.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("rendered template that does not exist")
	}
}
func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}
