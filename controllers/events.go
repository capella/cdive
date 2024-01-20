package controllers

import (
	"net/http"

	"github.com/capella/cdive/models"
)

func (c *controllers) Events(w http.ResponseWriter, r *http.Request) {
	var events []models.Events
	c.DB.Find(&events)
	c.renderTemplate(
		"events",
		nil,
		events,
	)(w, r)
}

func (c *controllers) Event(w http.ResponseWriter, r *http.Request) {
	c.renderTemplate("event", nil, nil)(w, r)
}
