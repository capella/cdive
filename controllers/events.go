package controllers

import (
	"net/http"
	"strconv"

	"github.com/capella/cdive/models"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func (c *controllers) Events(w http.ResponseWriter, r *http.Request) {
	var events []models.Event
	c.DB.Order("id desc, start desc").Where("deleted_at is NULL").Find(&events)
	c.renderTemplate(
		"events",
		nil,
		events,
	)(w, r)
}

func (c *controllers) Event(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	c.crud(&event, "event", w, r)
}

func (c *controllers) crud(data interface{}, name string, w http.ResponseWriter, r *http.Request) {
	formErrors := []string{}

	err := r.ParseForm()
	if err != nil {
		logrus.Error(err)
		return
	}

	vars := mux.Vars(r)
	pathID, perr := strconv.Atoi(vars["id"])

	err = decoder.Decode(data, r.PostForm)
	if err != nil {
		logrus.Error(err)
		formErrors = append(formErrors, err.Error())
	}

	if r.Method == http.MethodPost {
		result := c.DB.Save(data)
		if result.Error != nil {
			logrus.Error(err)
			formErrors = append(formErrors, err.Error())
		}
	} else if r.Method == http.MethodDelete {
		if perr == nil {
			c.DB.Delete(data, pathID)
		}
	} else {
		if perr == nil {
			c.DB.Find(data, pathID)
		}
	}

	c.renderTemplate(name, formErrors, data)(w, r)
}
