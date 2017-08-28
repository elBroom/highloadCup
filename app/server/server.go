package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/elBroom/highloadCup/app/handler"
	"github.com/elBroom/highloadCup/app/router"
	"github.com/elBroom/highloadCup/app/workers"
)

func init() {
	workers.Wp.Run()
}

func RunHTTPServer(addr string) error {
	_router := mux.NewRouter()
	_router.NotFoundHandler = http.HandlerFunc(handler.NotFound)

	router.Routing(_router)
	return http.ListenAndServe(addr, _router)
}
