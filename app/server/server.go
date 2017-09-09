package server

import (
	"time"

	"github.com/elBroom/highloadCup/app/router"
	"github.com/valyala/fasthttp"
)

func RunHTTPServer(addr string) error {
	return fasthttp.ListenAndServe(addr,
		fasthttp.TimeoutHandler(router.RequestHandler(), 2*time.Second, "timeout"))
}
