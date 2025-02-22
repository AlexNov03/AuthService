package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type SessionDelivery interface {
	Login(w http.ResponseWriter, r *http.Request)
	SignUp(w http.ResponseWriter, r *http.Request)
	LogOut(w http.ResponseWriter, r *http.Request)
}

type TaskDelivery interface {
	AddTask(w http.ResponseWriter, r *http.Request)
	GetTasks(w http.ResponseWriter, r *http.Request)
}

type Middleware interface {
	RequireAuth(h http.Handler) http.Handler
}

func NewHandler(sd SessionDelivery, td TaskDelivery, m Middleware) http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/logout", sd.LogOut).Methods(http.MethodGet)
	router.HandleFunc("/login", sd.Login).Methods(http.MethodPost)
	router.HandleFunc("/signup", sd.SignUp).Methods(http.MethodPost)
	router.Handle("/task", m.RequireAuth(http.HandlerFunc(td.AddTask))).Methods(http.MethodPost)
	router.Handle("/tasks", m.RequireAuth(http.HandlerFunc(td.GetTasks))).Methods(http.MethodGet)
	return router
}
