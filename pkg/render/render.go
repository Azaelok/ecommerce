package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Azaelok/ecommerce/pkg/config"
	"github.com/Azaelok/ecommerce/pkg/models"
)

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {

	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {

	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("no puedo obtener el template, desde template")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("algo", err)
	}

}

func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {

		name := filepath.Base(page)
		//log.Println("name", name)

		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		//log.Println("ts", ts)

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, nil
		}
		//log.Println("matches", matches)

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
			//log.Println("ts_2", ts)
		}

		myCache[name] = ts

		//log.Println(myCache)
	}

	return myCache, nil
}
