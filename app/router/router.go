package router

import (
	"github.com/elBroom/highloadCup/app/handler"
	"github.com/gorilla/mux"
)

func Routing(router *mux.Router) {
	router.HandleFunc("/users/{id}", handler.GetUserEndpoint).Methods("GET")
	router.HandleFunc("/users/{id}/visits", handler.VisitUserEndpoint).Methods("GET")
	router.HandleFunc("/users/new", handler.CreateUserEndpoint).Methods("POST")
	router.HandleFunc("/users/{id}", handler.UpdateUserEndpoint).Methods("POST")

	router.HandleFunc("/locations/{id}", handler.GetLocationEndpoint).Methods("GET")
	router.HandleFunc("/locations/{id}/avg", handler.GetLocatioAvgnEndpoint).Methods("GET")
	router.HandleFunc("/locations/new", handler.CreateLocationEndpoint).Methods("POST")
	router.HandleFunc("/locations/{id}", handler.UpdateLocationEndpoint).Methods("POST")

	router.HandleFunc("/visits/{id}", handler.GetVisitEndpoint).Methods("GET")
	router.HandleFunc("/visits/new", handler.CreateVisitEndpoint).Methods("POST")
	router.HandleFunc("/visits/{id}", handler.UpdateVisitEndpoint).Methods("POST")
}
