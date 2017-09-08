package handler

import (
	"net/http"

	"github.com/elBroom/highloadCup/app/model"
	"github.com/elBroom/highloadCup/app/schema"
	"github.com/elBroom/highloadCup/app/storage"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
)

func GetUserEndpoint(ctx *fasthttp.RequestCtx, id uint32) {
	user, ok := storage.DataStorage.User.Get(id)
	if !ok {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}

	writeObj(ctx, user)
}

func VisitUserEndpoint(ctx *fasthttp.RequestCtx, id uint32) {
	params := ctx.QueryArgs()

	//  Parse fromDate parameter
	var fromDate *int64
	if params.Has("fromDate") {
		if _, err := params.GetUint("fromDate"); err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		tmp, _ := params.GetUint("fromDate")
		tmp2 := int64(tmp)
		fromDate = &tmp2
	}

	//  Parse toDate parameter
	var toDate *int64
	if params.Has("toDate") {
		if _, err := params.GetUint("toDate"); err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		tmp, _ := params.GetUint("toDate")
		tmp2 := int64(tmp)
		toDate = &tmp2
	}

	//  Parse country parameter
	var country *string
	if params.Has("country") {
		tmp := string(params.Peek("country"))
		if tmp == "" {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		country = &tmp
	}

	//  Parse toDistance parameter
	var toDistance *uint32
	if params.Has("toDistance") {
		if _, err := params.GetUint("toDistance"); err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		tmp, _ := params.GetUint("toDistance")
		tmp2 := uint32(tmp)
		toDistance = &tmp2
	}

	var resp (schema.ResponceUserVisits)
	resp.Visits = []*schema.ResponceUserVisit{}
	ok := storage.DataStorage.VisitList.GetByUser(id, fromDate, toDate, country, toDistance, &resp.Visits)
	if !ok {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}

	writeObj(ctx, resp)
}

func UpdateUserEndpoint(ctx *fasthttp.RequestCtx, id uint32) {
	bytes := ctx.PostBody()
	ok := CheckNull(bytes)
	if ok {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	var user model.User
	_ = easyjson.Unmarshal(bytes, &user)

	err := storage.DataStorage.User.Update(id, &user)
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
	_ = easyjson.Unmarshal(ctx.PostBody(), &user)

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
