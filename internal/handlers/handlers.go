package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Azaelok/ecommerce/internal/config"
	"github.com/Azaelok/ecommerce/internal/driver"
	"github.com/Azaelok/ecommerce/internal/forms"
	"github.com/Azaelok/ecommerce/internal/helpers"
	"github.com/Azaelok/ecommerce/internal/models"
	"github.com/Azaelok/ecommerce/internal/render"
	"github.com/Azaelok/ecommerce/internal/repository"
	"github.com/Azaelok/ecommerce/internal/repository/dbrepo"
	"github.com/go-chi/chi/v5"
)

var Repo *Reposository

type Reposository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

func NewRepo(a *config.AppConfig, db *driver.DB) *Reposository {
	return &Reposository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewTestRepo(a *config.AppConfig) *Reposository {
	return &Reposository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

func NewHandlers(r *Reposository) {
	Repo = r
}

// Pagina Inicio
func (m *Reposository) Home(w http.ResponseWriter, r *http.Request) {

	//Se comento en el manejo de errores
	//remoteIP := r.RemoteAddr
	//m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// Página About
func (m *Reposository) About(w http.ResponseWriter, r *http.Request) {

	//Se comento en el manejo de errores
	//stringMap := make(map[string]string)
	//stringMap["test"] = "Hola AZAEL"

	//remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	//stringMap["remote_ip"] = remoteIP

	//render.Template(w, r, "about.page.tmpl", &models.TemplateData{
	//	StringMap: stringMap,
	//})

	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Pagina Generals
func (m *Reposository) Generals(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Pagina Majors
func (m *Reposository) Majors(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Pagina Search
func (m *Reposository) Availability(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

func (m *Reposository) PostAvailability(w http.ResponseWriter, r *http.Request) {

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	//for _, i := range rooms {
	//	m.App.InfoLog.Println("ROOM", i.ID, i.RoomName)
	//}

	if len(rooms) == 0 {
		//m.App.InfoLog.Println("No avail")
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})

	//w.Write([]byte(fmt.Sprintf("start date is %s end date is %s", start, end)))
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// AvailabilityJson
func (m *Reposository) AvailabilityJson(w http.ResponseWriter, r *http.Request) {

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, _ := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		//log.Println(err)
		helpers.ServerError(w, err)
		return
	}

	//log.Println(string(out))

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Muestra la lista de habitaciones disponibles
func (m *Reposository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Get(r.Context(), "reservation")

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}

	res.RoomId = roomID

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

// Toma los parametros de la url y construye una variable de sesion
func (m *Reposository) BookRoom(w http.ResponseWriter, r *http.Request) {
	// id, s, e
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	var res models.Reservation

	room, err := m.DB.GetRoomID(roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.Room.RoomName = room.RoomName
	res.RoomId = roomID
	res.StartDate = startDate
	res.EndDate = endDate

	m.App.Session.Put(r.Context(), "reservation", res)

	log.Println(roomID, startDate, endDate)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

// Pagina Contact
func (m *Reposository) Contact(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// Pagina Make-reservation
func (m *Reposository) Reservation(w http.ResponseWriter, r *http.Request) {

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		//helpers.ServerError(w, errors.New("cannot get reservation from session!"))
		return
	}

	room, err := m.DB.GetRoomID(res.RoomId)
	if err != nil {
		//helpers.ServerError(w, err)
		//return
		m.App.Session.Put(r.Context(), "error", "Can't find room!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		//helpers.ServerError(w, errors.New("cannot get reservation from session!"))
		return
	}

	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// Pagina Make-reservation
func (m *Reposository) PostReservation(w http.ResponseWriter, r *http.Request) {

	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("can't get from session"))
		return
	}

	err := r.ParseForm()
	//err = errors.New("This is an error message")
	if err != nil {
		//log.Println(err)
		helpers.ServerError(w, err)
		return
	}

	// sd := r.Form.Get("start_date")
	// ed := r.Form.Get("end_date")

	// layout := "2006-01-02"
	// startDate, err := time.Parse(layout, sd)
	// if err != nil {
	// 	helpers.ServerError(w, err)
	// 	return
	// }

	// endDate, err := time.Parse(layout, ed)
	// if err != nil {
	// 	helpers.ServerError(w, err)
	// 	return
	// }

	// roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	// if err != nil {
	// 	helpers.ServerError(w, err)
	// 	return
	// }

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	// reservation := models.Reservation{
	// 	FirstName: r.Form.Get("first_name"),
	// 	LastName:  r.Form.Get("last_name"),
	// 	Email:     r.Form.Get("email"),
	// 	Phone:     r.Form.Get("phone"),
	// 	StartDate: startDate,
	// 	EndDate:   endDate,
	// 	RoomId:    roomID,
	// }

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomId:        reservation.RoomId,
		ReservationId: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

// Muestra una pagina con el resumen de la reserva
func (m *Reposository) ReservationSummary(w http.ResponseWriter, r *http.Request) {

	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		//log.Println("cannot get item from session")
		m.App.ErrorLog.Println("cannot get item from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from sessión")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Eliminamos los datos de nuestra reserva
	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}
