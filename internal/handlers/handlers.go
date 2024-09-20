package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Azaelok/ecommerce/internal/config"
	"github.com/Azaelok/ecommerce/internal/models"
	"github.com/Azaelok/ecommerce/internal/render"
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

	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

// Página About
func (m *Reposository) About(w http.ResponseWriter, r *http.Request) {

	stringMap := make(map[string]string)
	stringMap["test"] = "Hola AZAEL"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Pagina Generals
func (m *Reposository) Generals(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{
		//		StringMap: stringMap,
	})
}

// Pagina Majors
func (m *Reposository) Majors(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{
		//		StringMap: stringMap,
	})
}

// Pagina Search
func (m *Reposository) Search(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{
		//		StringMap: stringMap,
	})
}

func (m *Reposository) PostSearch(w http.ResponseWriter, r *http.Request) {

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("start date is %s end date is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJson
func (m *Reposository) AvailabilityJson(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      false,
		Message: "Disponible",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Println(err)
	}

	log.Println(string(out))

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Pagina Contact
func (m *Reposository) Contact(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{
		//		StringMap: stringMap,
	})
}

// Pagina Make-reservation
func (m *Reposository) Make(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		//		StringMap: stringMap,
	})
}
