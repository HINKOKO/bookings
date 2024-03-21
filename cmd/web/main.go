package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/HINKOKO/bookings/pkg/config"
	"github.com/HINKOKO/bookings/pkg/handlers"
	"github.com/HINKOKO/bookings/pkg/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	// Change this to true when in production
	app.InProduction = false

	// variable shadowing
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true                  // Should cookie persist after window is closed ? true
	session.Cookie.SameSite = http.SameSiteLaxMode //
	session.Cookie.Secure = app.InProduction       // Cookie crypted https instead of http -> true for production, false for dev

	app.Session = session
	// session.Cookie.Name = "session_id"
	// session.Cookie.Domain = "example.com"
	// session.Cookie.HttpOnly = true
	// session.Cookie.Path = "/example/"

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	// Now handled by 'routes.go' and package pat
	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println(fmt.Sprintf("starting app on port %s", portNumber))
	// _ = http.ListenAndServe(portNumber, nil)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
