package handlers

import (
	"net/http"

	"github.com/HINKOKO/bookings/pkg/config"
	"github.com/HINKOKO/bookings/pkg/models"
	"github.com/HINKOKO/bookings/pkg/render"
	// "github.com/HINKOKO/go-course/pkg/render"
)

// Repo -> the repository used by the handlers
var Repo *Repository

// repository type
type Repository struct {
	App *config.AppConfig
}

// Creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// @NEwHandlers: sets the repository for the handler
func NewHandlers(r *Repository) {
	Repo = r
}

// Handler for Home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIp)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

// Handler for About page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "hello , again and again"

	// We can now access the 'session'
	// m.App.Session.Cookie = blabla

	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp

	// send the data to the template
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
