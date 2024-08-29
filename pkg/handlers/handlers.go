package handlers

import (
	"net/http"

	"github.com/Azaelok/ecommerce/pkg/config"
	"github.com/Azaelok/ecommerce/pkg/models"
	"github.com/Azaelok/ecommerce/pkg/render"
)

var Repo *Reposository

type Reposository struct {
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Reposository {

	return &Reposository{
		App: a,
	}
}

func NewHandlers(r *Reposository) {
	Repo = r
}

// Pagina Inicio
func (m *Reposository) Home(w http.ResponseWriter, r *http.Request) {

	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

// Página About
func (m *Reposository) About(w http.ResponseWriter, r *http.Request) {

	stringMap := make(map[string]string)
	stringMap["test"] = "Hola AZAEL"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
