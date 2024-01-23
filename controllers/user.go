package controllers

import (
	"net/http"
)

func (c *controllers) User(w http.ResponseWriter, r *http.Request) {
	user := getContextUser(r.Context())
	c.renderTemplate(
		"user",
		nil,
		user,
	)(w, r)
}
