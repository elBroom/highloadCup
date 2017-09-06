package router

import (
	"strconv"

	"github.com/elBroom/highloadCup/app/handler"
	"github.com/golang/glog"
	"github.com/valyala/fasthttp"
)

var Phase = 1
var isLastGet = true

func RequestHandler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		path := ctx.Path()
		if ctx.IsGet() {
			if !isLastGet {
				Phase++
				isLastGet = true
				glog.Infof("Phase: %d", Phase)
			}
			if matchGetUser(path) > 0 {
				// /users/:id
				id, err := parseID(path[7:])
				if err == nil {
					handler.GetUserEndpoint(ctx, id)
					return
				}
			} else if matchGetLocation(path) > 0 {
				// /locations/:id
				id, err := parseID(path[11:])
				if err == nil {
					handler.GetLocationEndpoint(ctx, id)
					return
				}
			} else if matchGetVisit(path) > 0 {
				// /visits/:id
				id, err := parseID(path[8:])
				if err == nil {
					handler.GetVisitEndpoint(ctx, id)
					return
				}
			} else if matchLocationAvg(path) > 0 {
				// /locations/:id/avg
				id, err := parseID(path[11 : len(path)-4])
				if err == nil {
					handler.GetLocatioAvgEndpoint(ctx, id)
					return
				}
			} else if matchVisitUser(path) > 0 {
				// /users/:id/visits
				id, err := parseID(path[7 : len(path)-7])
				if err == nil {
					handler.VisitUserEndpoint(ctx, id)
					return
				}
			}
		} else if ctx.IsPost() {
			if isLastGet {
				Phase++
				isLastGet = false
				glog.Infof("Phase: %d", Phase)
			}
			if matchCreateUser(path) > 0 {
				// /users/new
				handler.CreateUserEndpoint(ctx)
				return
			} else if matchCreateLocation(path) > 0 {
				// /loactions/new
				handler.CreateLocationEndpoint(ctx)
				return
			} else if matchCreateVisit(path) > 0 {
				// /visits/new
				handler.CreateVisitEndpoint(ctx)
				return
			} else if matchUpdateUser(path) > 0 {
				// /users/:id
				id, err := parseID(path[7:])
				if err == nil {
					handler.UpdateUserEndpoint(ctx, id)
					return
				}
			} else if matchUpdateLocation(path) > 0 {
				// /locations/:id
				id, err := parseID(path[11:])
				if err == nil {
					handler.UpdateLocationEndpoint(ctx, id)
					return
				}
			} else if matchUpdateVisit(path) > 0 {
				// /visits/:id
				id, err := parseID(path[8:])
				if err == nil {
					handler.UpdateVisitEndpoint(ctx, id)
					return
				}
			}
		}
		ctx.SetStatusCode(404)
	}
}

func parseID(str []byte) (uint32, error) {
	id, err := strconv.ParseUint(string(str), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}
