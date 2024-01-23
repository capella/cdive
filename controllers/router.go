package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (c *controllers) Router() *mux.Router {
	router := mux.NewRouter()
	router.Use(c.SessionMiddleware)

	router.HandleFunc("/login", c.LoginPOST).Methods(http.MethodPost)
	router.HandleFunc("/login", c.Login)

	router.HandleFunc("/events/{category:.*}", c.Events)
	router.HandleFunc("/events", c.Events)

	router.Use(mux.CORSMethodMiddleware(router))

	logged := router.PathPrefix("/").Subrouter()
	logged.Use(c.Internal)
	logged.HandleFunc("/logout", c.Logout)
	logged.HandleFunc("/", c.Home)

	logged.HandleFunc("/event/{id}", c.Event)
	logged.HandleFunc("/event/", c.Event)

	logged.HandleFunc("/me", c.User)

	admin := router.PathPrefix("/").Subrouter()
	admin.Use(c.Admin)
	admin.HandleFunc("/admin", c.Home)
	admin.HandleFunc("/members", c.Members)
	admin.HandleFunc("/members/{id}", c.Member)

	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("."))),
	)

	return router
}
