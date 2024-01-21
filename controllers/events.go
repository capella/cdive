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
	c.crud(w, r)
}

func (c *controllers) crud(w http.ResponseWriter, r *http.Request) {
	formErrors := []string{}

	var event models.Events
	err := r.ParseForm()
	if err != nil {
		logrus.Error(err)
		formErrors = append(formErrors, err.Error())
		return
	}

	vars := mux.Vars(r)
	pathID, perr := strconv.Atoi(vars["id"])

	err = decoder.Decode(&event, r.PostForm)
	if err != nil {
		logrus.Error(err)
		formErrors = append(formErrors, err.Error())
	}

	if r.Method == http.MethodPost {
		result := c.DB.Save(&event)
		if result.Error != nil {
			logrus.Error(err)
			formErrors = append(formErrors, err.Error())
		}
	} else if r.Method == http.MethodDelete {
		if perr == nil {
			c.DB.Delete(&event, pathID)
		}
	} else {
		if perr == nil {
			c.DB.Find(&event, pathID)
		}
	}

	c.renderTemplate("event", formErrors, event)(w, r)
}
