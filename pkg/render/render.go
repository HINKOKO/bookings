package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/HINKOKO/bookings/pkg/config"
	"github.com/HINKOKO/bookings/pkg/models"
)

// var functions = template.FuncMap{}

var app *config.AppConfig

// Sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
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
		log.Fatal("Could not get requested template from template cache")
	}

	// For finner grained error management
	buf := new(bytes.Buffer)

	td = AddDefaultData(td)
	_ = t.Execute(buf, td)

	// Render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error: writing template to the browser", err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	// myCache := make(map[string]*template.Template)
	myCache := map[string]*template.Template{} // exactly the same as previous syntax

	// get all of the files name *.page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	// log.Println(pages)
	if err != nil {
		return myCache, err
	}

	// Range through the 'pages'.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
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
