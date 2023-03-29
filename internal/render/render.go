package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/psampath6/bookings/internal/config"

	//"github.com/psampath6/bookings/internal/handlers"
	"github.com/psampath6/bookings/internal/models"
)

var functions = template.FuncMap {

}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds data for all templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	return td
}
// RenderTemplate renders a templates using html/template
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
    var tc map[string]*template.Template
	if app.UseCache {
	    // get the template cache from the app config
	    tc = app.TemplateCache
	} else {
		// this is just used for testing, so that we rebuild
		// the cache on every request
		tc, _ = CreateTemplateCache()
	}

	// create a template
    //tc, err := CreateTemplateCache()
	//if err != nil {
	//	log.Fatal(err)
	//}

	// get requested template from cache
    t, ok := tc[tmpl]
	if !ok {
		//log.Fatal(err)
		log.Fatal("Could not get template from template cache")
	}
	buf := new(bytes.Buffer)
	td = AddDefaultData(td, r)
	_ = t.Execute(buf, td)
	
	// render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}
	return err
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	//myCache := make(map[string]*template.Template)
	myCache := map[string]*template.Template{}

	// get all of the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")

	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.tmpl

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
/*

// RenderTemplate renders a template
func RenderTemplateTest(w http.ResponseWriter, tmpl string) {
	parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl, "./templates/base.layout.tmpl")
	err := parsedTemplate.Execute(w, nil)
	if err != nil {
		fmt.Println("error parsing template:", err)
		return
	}
}
var tc = make(map[string]*template.Template)
func RenderTemplate(w http.ResponseWriter, t string) {
    var tmpl *template.Template
	var err error

	// check to see if we already have the template in our cache
	_, inMap := tc[t]
	if !inMap {
		// need to create the template
		log.Println("creating template and adding to cache")
		err = createTemplateCache(t)
		if err != nil {
			log.Println(err)
		}
	} else {
		// we have the template in the cache
		log.Println("using cached template")
	}

	tmpl = tc[t]

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

func createTemplateCache(t string) error {
	templates := []string{
		fmt.Sprintf("./templates/%s", t), 
		"./templates/base.layout.tmpl",
	}

	// parse the template
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}

	// add template to cache (map)
	tc[t] = tmpl

	return nil
}
*/