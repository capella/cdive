package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (c *controllers) Router() *mux.Router {
	router := mux.NewRouter()
	router.Use(c.SessionMiddleware)

	router.HandleFunc("/login", c.LoginPOST).Methods("POST")
	router.HandleFunc("/login", c.Login)

	router.Use(mux.CORSMethodMiddleware(router))

	logged := router.PathPrefix("/").Subrouter()
	logged.Use(c.Internal)
	logged.HandleFunc("/logout", c.Logout)
	logged.HandleFunc("/", c.Home)

	// TODO(capella): add controllers
	admin := router.PathPrefix("/").Subrouter()
	admin.Use(c.Admin)
	logged.HandleFunc("/admin", c.Home)

	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("."))),
	)

	return router
}
