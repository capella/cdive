package controllers

import (
	"net/http"
	"strconv"

	"github.com/capella/cdive/models"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func (c *controllers) Events(w http.ResponseWriter, r *http.Request) {
	var events []models.Events
	c.DB.Order("id desc, start desc").Where("deleted_at is NULL").Find(&events)
	c.renderTemplate(
		"events",
		nil,
		events,
	)(w, r)
}

func (c *controllers) Event(w http.ResponseWriter, r *http.Request) {
	var event models.Events
	err := r.ParseForm()
	if err != nil {
		logrus.Error(err)
		return
	}

	vars := mux.Vars(r)
	pathID, perr := strconv.Atoi(vars["id"])

	decoder.Decode(&event, r.PostForm)
	if r.Method == http.MethodPost {
		c.DB.Save(&event)
	} else if r.Method == http.MethodDelete {
		if perr == nil {
			c.DB.Delete(&event, pathID)
		}
	} else {
		if perr == nil {
			c.DB.Find(&event, pathID)
		}
	}
	c.renderTemplate("event", nil, event)(w, r)
}
