package router

import (
	"strconv"

	"github.com/elBroom/highloadCup/app"
	"github.com/elBroom/highloadCup/app/handler"
	"github.com/elBroom/highloadCup/app/storage"
	"github.com/valyala/fasthttp"
)

var isLastGet = true

func RequestHandler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		path := ctx.Path()
		if ctx.IsGet() {
			if !isLastGet {
				app.Phase++
				isLastGet = true
			}

			if matchGetUser(path) > 0 {
				// /users/:id
				id, err := parseID(path[7:])
				if id <= storage.CountUser && err == nil {
					handler.GetUserEndpoint(ctx, id)
					return
				}
			} else if matchGetLocation(path) > 0 {
				// /locations/:id
				id, err := parseID(path[11:])
				if id <= storage.CountLocation && err == nil {
					handler.GetLocationEndpoint(ctx, id)
					return
				}
			} else if matchGetVisit(path) > 0 {
				// /visits/:id
				id, err := parseID(path[8:])
				if id <= storage.CountVisit && err == nil {
					handler.GetVisitEndpoint(ctx, id)
					return
				}
			}

			resp, ok := app.MemoryCache.Get(ctx.URI())
			if ok {
				ctx.SetStatusCode(*resp.StatusCode())
				ctx.Response.Header.Set("Content-Type", "application/json")
				ctx.Response.Header.Set("Connection", "keep-alive")
				ctx.Response.Header.Set("Content-Length", strconv.Itoa(len(*resp.Body())))
				ctx.SetBody(*resp.Body())
				return
			}
			if matchLocationAvg(path) > 0 {
				// /locations/:id/avg
				id, err := parseID(path[11 : len(path)-4])
				if id <= storage.CountLocation && err == nil {
					handler.GetLocatioAvgEndpoint(ctx, id)
					app.MemoryCache.Set(ctx.URI(), &ctx.Response, app.LifeTime)
					return
				}
			} else if matchVisitUser(path) > 0 {
				// /users/:id/visits
				id, err := parseID(path[7 : len(path)-7])
				if id <= storage.CountUser && err == nil {
					handler.VisitUserEndpoint(ctx, id)
					app.MemoryCache.Set(ctx.URI(), &ctx.Response, app.LifeTime)
					return
				}
			}
		} else if ctx.IsPost() {
			if isLastGet {
				app.Phase++
				isLastGet = false
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
				if id <= storage.CountUser && err == nil {
					handler.UpdateUserEndpoint(ctx, id)
					return
				}
			} else if matchUpdateLocation(path) > 0 {
				// /locations/:id
				id, err := parseID(path[11:])
				if id <= storage.CountLocation && err == nil {
					handler.UpdateLocationEndpoint(ctx, id)
					return
				}
			} else if matchUpdateVisit(path) > 0 {
				// /visits/:id
				id, err := parseID(path[8:])
				if id <= storage.CountVisit && err == nil {
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
