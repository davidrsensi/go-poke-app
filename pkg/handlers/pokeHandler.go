package handlers

import (
	"go-poke-app/pkg/models"
	"go-poke-app/pkg/utils"
	"log"
	"net/http"
	"time"
)

// Pokes is a http.Handler
type PokeHandler struct {
	l *log.Logger
}

// NewPokes creates a products handler with the given logger
func NewPokes(l *log.Logger) *PokeHandler {
	return &PokeHandler{l}
}

// GetPokes returns the pokes from mongoDB
func (ph *PokeHandler) GetPokes(rw http.ResponseWriter, r *http.Request) {
	ph.l.Println("Handle GetPokes")

	// fetch the pokes from the datastore
	lp := models.GetPokes()

	utils.ExecuteTemplate(rw, "index.html", lp)

}

// add poke to DB and send email
func (ph *PokeHandler) PostPokes(rw http.ResponseWriter, r *http.Request) {
	ph.l.Println("Handle PostPokes")

	r.ParseForm()

	postPoke := models.Poke{
		SenderName:    r.FormValue("sender_name"),
		ReceiverName:  r.FormValue("receiver_name"),
		ReceiverEmail: r.FormValue("receiver_email"),
		Message:       r.FormValue("message"),
		DateSent:      time.Now(),
	}

	if r.FormValue("public") == "1" {
		postPoke.IsPublic = true
	} else {
		postPoke.IsPublic = false
	}

	// ############# Should probably send email first and if it fails, do not create the DB entry.
	if err := postPoke.Create(); err == nil {
		postPoke.SendEmail()
	}

	http.Redirect(rw, r, "/", 302)

}
