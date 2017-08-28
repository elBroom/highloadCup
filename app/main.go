package main

import (
	"log"
	"net/http"

	"github.com/elBroom/highloadCup/app/handler"
	"github.com/elBroom/highloadCup/app/importer"
	"github.com/elBroom/highloadCup/app/router"
	"github.com/gorilla/mux"
)

func main() {
	_router := mux.NewRouter()
	_router.NotFoundHandler = http.HandlerFunc(handler.NotFound)

	err := importer.ImportDataFromZip()
	if err != nil {
		log.Fatal(err)
	}

	router.Routing(_router)
	srv := &http.Server{Addr: ":8000", Handler: _router}
	log.Println("Start server")
	log.Fatal(srv.ListenAndServe())
}
