package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/HINKOKO/bookings/internal/config"
	"github.com/HINKOKO/bookings/internal/driver"
	"github.com/HINKOKO/bookings/internal/handlers"
	"github.com/HINKOKO/bookings/internal/helpers"
	"github.com/HINKOKO/bookings/internal/models"
	"github.com/HINKOKO/bookings/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	// Better place here for defering database closing
	defer db.SQL.Close()

	defer close(app.MailChan)
	listenForMail()

	// test to send email
	msg := models.MailData{
		To:      "john@doe.ca",
		From:    "me@here.com",
		Subject: "Some",
		Content: "content",
	}
	app.MailChan <- msg

	fmt.Printf("starting app on port %s\n", portNumber)
	// _ = http.ListenAndServe(portNumber, nil)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// What am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// Change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// variable shadowing
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true                  // Should cookie persist after window is closed ? true
	session.Cookie.SameSite = http.SameSiteLaxMode //
	session.Cookie.Secure = app.InProduction       // Cookie crypted https instead of http -> true for production, false for dev

	app.Session = session

	// Connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 database=bookings user=postgres password=")
	if err != nil {
		log.Fatal("Unable to connect database.Dying...")
	}
	log.Println("Connected to database dude.")

	// defer db.SQL.Close() If we put this here  -- run finish 'running' and BAM you close the DB, not good

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, err
}
