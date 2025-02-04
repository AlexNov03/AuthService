package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Delivery interface {
	Login(w http.ResponseWriter, r *http.Request)
	SignUp(w http.ResponseWriter, r *http.Request)
	LogOut(w http.ResponseWriter, r *http.Request)
}

func NewHandler(delivery Delivery) http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/logout", delivery.LogOut).Methods(http.MethodGet)
	router.HandleFunc("/login", delivery.Login).Methods(http.MethodPost)
	router.HandleFunc("/signup", delivery.SignUp).Methods(http.MethodPost)
	return router
}
