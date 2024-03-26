package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/HINKOKO/bookings/internal/config"
	"github.com/HINKOKO/bookings/internal/forms"
	"github.com/HINKOKO/bookings/internal/models"
	"github.com/HINKOKO/bookings/internal/render"
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

	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
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
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// handler for Contacxt page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// handler for Reservation page
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyRes models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyRes

	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation - handles the posting of reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm() // PArse the form data
	if err != nil {
		log.Println(err)
		return
	}
	// Populate a reservation variable with the form data
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}
	// This might be bad data ...
	// PostForm contains the parsed form data from PATCH, POST
	// or PUT body parameters.
	//
	// This field is only available after ParseForm is called.
	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 4)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	// take users to reservation summary page
	// take advantage of the SESSION
	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// handler for generals quarters room
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// handler for major suite room
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// handler for major suite room
func (m *Repository) Available(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// handler for major suite room
func (m *Repository) PostAvailable(w http.ResponseWriter, r *http.Request) {
	// Capture information from the forms of Html page
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// handles request for available dates and send JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "available!",
	}

	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		log.Println(err)
	}
	log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	// type assertion trick with '.' connection
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		log.Println("cannot get item from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.RenderTemplate(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
