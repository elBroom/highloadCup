package main

import (
	"log"

	"runtime"

	"github.com/elBroom/highloadCup/app/importer"
	"github.com/elBroom/highloadCup/app/server"
)

func main() {
	err := importer.ImportDataFromZip()
	if err != nil {
		log.Fatal(err)
	}

	runtime.GC()
	port := ":80"
	log.Printf("Start server on %s", port)
	log.Fatal(server.RunHTTPServer(port))
}
