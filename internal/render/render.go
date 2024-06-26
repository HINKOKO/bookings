package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/HINKOKO/bookings/internal/config"
	"github.com/HINKOKO/bookings/internal/models"
	"github.com/justinas/nosurf"
)

// 'functions' used to pass functions available to go templates
var functions = template.FuncMap{
	"humanDate":  HumanDate,
	"formatDate": FormatDate,
	"iterate":    Iterate,
	"add":        Add,
}
var app *config.AppConfig
var pathToTemplates = "./templates"

func Add(a, b int) int {
	return a + b
}

// Iterate - returns a slice of ints, starting at 1 to count
func Iterate(count int) []int {
	var i int
	var items []int

	for i = 0; i < count; i++ {
		items = append(items, i)
	}
	return items
}

// NewRenderer - sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// HumanDate - Handles database ugly date format and convert to nicely readable date
func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDate - formats a date with given string as param
func FormatDate(t time.Time, f string) string {
	return t.Format(f)
}

// AddDefaultData - adds data for all the templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	}

	return td
}

// Template - renders templates using html/templates
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template
	if app.UseCache {
		// get the template from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// Get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		// log.Fatal("Could not get requested template from template cache")
		return errors.New("can't get template from cache")
	}

	// For finner grained error management
	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)
	_ = t.Execute(buf, td)

	// Render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error: writing template to the browser", err)
		return err
	}

	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	// myCache := make(map[string]*template.Template)
	myCache := map[string]*template.Template{} // exactly the same as previous syntax

	// get all of the files name *.page.tmpl from ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	// log.Println(pages)
	if err != nil {
		return myCache, err
	}

	// Range through the 'pages'.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}

// ====== BASIC CACHING SYSTEM, like DO IT YOURSELF HOMEMADE LEROY MERLIN ======
// 'template.Template' represents a parsed template, struct type that holds
// parse repres of a template and provides methods for executing that template against data
// var tc = make(map[string]*template.Template)

// func RenderTemplate(w http.ResponseWriter, t string) {
// 	var tmpl *template.Template
// 	var err error

// 	// Check to see if we already have the template in our cache
// 	_, inMap := tc[t]
// 	if !inMap {
// 		log.Println("creating template and adding to cache")

// 		// Need to create the tmpl
// 		err = createTemplateCache(t)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	} else {
// 		// We have template in the cache
// 		log.Println("using cached template")

// 	}

// 	tmpl = tc[t]
// 	err = tmpl.Execute(w, nil)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

// func createTemplateCache(t string) error {
// 	templates := []string{
// 		fmt.Sprintf("./templates/%s", t),
// 		"./templates/base.layout.tmpl",
// 	}
// 	// parse the templates, ParseFiles returns pointer to template.Template instance
// 	tmpl, err := template.ParseFiles(templates...) // '...' take each strings from 'templates' and treats them as individuals
// 	if err != nil {
// 		return err
// 	}
// 	tc[t] = tmpl

// 	return nil
// }
