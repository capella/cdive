package controllers

import (
	"net/http"

	"github.com/capella/cdive/models"
)

func (c *controllers) Members(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	c.DB.Order("name desc").Find(&users)
	c.renderTemplate(
		"members",
		nil,
		users,
	)(w, r)
}

func (c *controllers) Member(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	c.crud(&user, "member", w, r)
}
