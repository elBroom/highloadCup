package router

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/elBroom/highloadCup/app/handler"
)

func Routing() *fasthttprouter.Router {
	router := fasthttprouter.New()

	router.GET("/users/:id", handler.GetUserEndpoint)
	router.GET("/users/:id/visits", handler.VisitUserEndpoint)
	router.POST("/users/new", handler.CreateUserEndpoint)

	router.GET("/locations/:id", handler.GetLocationEndpoint)
	router.GET("/locations/:id/avg", handler.GetLocatioAvgnEndpoint)
	router.POST("/locations/new", handler.CreateLocationEndpoint)

	router.GET("/visits/:id", handler.GetVisitEndpoint)
	router.POST("/visits/new", handler.CreateVisitEndpoint)

	router2 := fasthttprouter.New()
	router.MethodNotAllowed = router2.Handler

	router2.POST("/users/:id", handler.UpdateUserEndpoint)
	router2.POST("/locations/:id", handler.UpdateLocationEndpoint)
	router2.POST("/visits/:id", handler.UpdateVisitEndpoint)
	return router
}
