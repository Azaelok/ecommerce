package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Azaelok/ecommerce/internal/config"
	"github.com/Azaelok/ecommerce/internal/handlers"
	"github.com/Azaelok/ecommerce/internal/render"
	"github.com/alexedwards/scs/v2"
)

// Uso de constantes
const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("Start application on port #{portNumber}"))

	srv := http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() error {

	// Cambiar a true cuando este en produccion
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("No se puede crear la plantilla en cache")
		return err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	return nil
}
