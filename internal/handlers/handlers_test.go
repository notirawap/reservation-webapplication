package handlers

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/notirawap/doctor-registration/internal/driver"
)

func TestNewRepo(t *testing.T) {
	var db driver.DB
	testRepo := NewRepo(&app, &db)

	if reflect.TypeOf(testRepo).String() != "*handlers.Repository" {
		t.Errorf("Did not get correct type")
	}
}

// Run table test to test all GET handlers
var theTests = []struct {
	name               string // test name
	url                string // path
	method             string // GET or POST
	expectedStatusCode int    // status code responses whether the test has passed
}{
	// GET test data entries
	{"home", "/", "GET", http.StatusOK},
	{"doctor-form", "/doctor/form", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"user-login", "/user/login", "GET", http.StatusOK},
	{"user-logout", "/user/logout", "GET", http.StatusOK},
}

func TestGETHandlers(t *testing.T) {
	// Create test server that can be posted to and returns status code
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	// Run Table test that makes GET request
	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			// Error if the status code is not matched
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
