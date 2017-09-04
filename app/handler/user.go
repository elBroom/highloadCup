package handler

import (
	"net/http"

	"sort"

	"github.com/elBroom/highloadCup/app/model"
	"github.com/elBroom/highloadCup/app/schema"
	"github.com/elBroom/highloadCup/app/storage"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/valyala/fasthttp"
)

func GetUserEndpoint(ctx *fasthttp.RequestCtx) {
	id, err := parseID(ctx)
	if err != nil {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}

	user, ok := storage.DataStorage.User.Get(id)
	if !ok {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}

	writeObj(ctx, user)
}

func VisitUserEndpoint(ctx *fasthttp.RequestCtx) {
	id, err := parseID(ctx)
	if err != nil {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}

	params := ctx.QueryArgs()

	//  Parse fromDate parameter
	var fromDate int64
	if params.Has("fromDate") {
		if _, err := params.GetUint("fromDate"); err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		tmp, _ := params.GetUint("fromDate")
		fromDate = int64(tmp)
	}

	//  Parse toDate parameter
	var toDate int64
	if params.Has("toDate") {
		if _, err := params.GetUint("toDate"); err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		tmp, _ := params.GetUint("toDate")
		toDate = int64(tmp)
	}

	//  Parse country parameter
	var country string
	if params.Has("country") {
		country = string(params.Peek("country"))
		if country == "" {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
	}

	//  Parse toDistance parameter
	var toDistance uint32
	if params.Has("toDistance") {
		if _, err := params.GetUint("toDistance"); err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		tmp, _ := params.GetUint("toDistance")
		toDistance = uint32(tmp)
	}

	visits, ok := storage.DataStorage.VisitList.GetByUser(id, storage.DataStorage)
	if !ok {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}

	var resp (schema.ResponceUserVisits)
	resp.Visits = []*schema.ResponceUserVisit{}
	for _, visit := range visits {
		if visit.LocationID == nil {
			continue
		}
		location, ok := storage.DataStorage.Location.Get(*visit.LocationID)
		if !ok {
			continue
		}
		if params.Has("fromDate") && fromDate > (*visit.VisitedAt) {
			continue
		}
		if params.Has("toDate") && toDate < (*visit.VisitedAt) {
			continue
		}
		if params.Has("country") && country != (*location.Country) {
			continue
		}
		if params.Has("toDistance") && toDistance <= (*location.Distance) {
			continue
		}
		var item schema.ResponceUserVisit
		item.Mark = visit.Mark
		item.Visited_at = visit.VisitedAt
		item.Place = location.Place
		resp.Visits = append(resp.Visits, &item)
	}

	sort.Sort(&resp)
	writeObj(ctx, resp)
}

func UpdateUserEndpoint(ctx *fasthttp.RequestCtx) {
	id, err := parseID(ctx)
	if err != nil {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}

	bytes := ctx.PostBody()
	ok := CheckNull(bytes)
	if ok {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	var user model.User
	_ = ffjson.Unmarshal(bytes, &user)

	err = storage.DataStorage.User.Update(id, &user)
	if err != nil {
		if err == storage.ErrDoesNotExist {
			ctx.SetStatusCode(http.StatusNotFound)
		} else {
			ctx.SetStatusCode(http.StatusBadRequest)
		}
		return
	}

	writeStr(ctx, "{}")
}

func CreateUserEndpoint(ctx *fasthttp.RequestCtx) {
	var user model.User
	_ = ffjson.Unmarshal(ctx.PostBody(), &user)

	err := storage.DataStorage.User.Add(&user)
	if err != nil {
		if err == storage.ErrDoesNotExist {
			ctx.SetStatusCode(http.StatusNotFound)
		} else {
			ctx.SetStatusCode(http.StatusBadRequest)
		}
		return
	}

	writeStr(ctx, "{}")
}
