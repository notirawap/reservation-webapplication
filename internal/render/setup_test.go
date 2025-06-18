package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/notirawap/doctor-registration/internal/config"
	"github.com/notirawap/doctor-registration/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig

// Copy from main.go
func TestMain(m *testing.M) {
	gob.Register(models.User{})
	gob.Register(models.Doctor{})

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}

// Create a type that satisfied the ResponseWriter interface
type myWriter struct{}

// Need Header(), Write(), WriterHeader() methods
func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}
func (tw *myWriter) WriteHeader(i int) {

}
func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
