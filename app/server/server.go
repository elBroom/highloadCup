package server

import (
	"github.com/elBroom/highloadCup/app/router"
	"github.com/valyala/fasthttp"
)

func RunHTTPServer(addr string) error {
	//router := router.Routing()

	return fasthttp.ListenAndServe(addr, router.RequestHandler())
}
