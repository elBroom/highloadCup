package handler

import (
	"encoding/json"
	"net/http"

	"sort"

	"io/ioutil"

	"math"
	"net/url"
	"strconv"

	"github.com/elBroom/highloadCup/app/model"
	"github.com/elBroom/highloadCup/app/schema"
	"github.com/elBroom/highloadCup/app/storage"
	"github.com/elBroom/highloadCup/app/workers"
)

func GetUserEndpoint(w http.ResponseWriter, req *http.Request) {
	_, err := workers.Wp.AddTaskSyncTimed(func() interface{} {
		id, err := parseID(req)
		if err != nil {

			http.Error(w, "", http.StatusNotFound)
			return nil
		}

		user, ok := storage.DataStorage.User.Get(id)
		if !ok {
			http.Error(w, "", http.StatusNotFound)
			return nil
		}

		json.NewEncoder(w).Encode(user)
		return nil
	}, workers.TimeOut)

	checkTimeout(w, err)
}

func VisitUserEndpoint(w http.ResponseWriter, req *http.Request) {
	_, err := workers.Wp.AddTaskSyncTimed(func() interface{} {
		id, err := parseID(req)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return nil
		}

		params, err := url.ParseQuery(req.URL.RawQuery)
		//  Parse fromDate parameter
		fromDateStr := params.Get("fromDate")
		fromDate := int64(0)

		if fromDateStr != "" {
			if fromDate, err = strconv.ParseInt(fromDateStr, 10, 64); err != nil {
				http.Error(w, "", http.StatusBadRequest)
				return nil
			}
		}

		//  Parse toDate parameter
		toDateStr := params.Get("toDate")
		toDate := int64(math.MaxInt64)

		if toDateStr != "" {
			if toDate, err = strconv.ParseInt(toDateStr, 10, 64); err != nil {
				http.Error(w, "", http.StatusBadRequest)
				return nil
			}
		}

		//  Parse country parameter
		country := params.Get("country")

		//  Parse toDistance parameter
		toDistStr := params.Get("toDistance")
		toDistance := int64(math.MaxInt32)

		if toDistStr != "" {
			if toDistance, err = strconv.ParseInt(toDistStr, 10, 64); err != nil {
				http.Error(w, "", http.StatusBadRequest)
				return nil
			}
		}

		visits, ok := storage.DataStorage.VisitList.GetByUser(id, storage.DataStorage)
		if !ok {
			http.Error(w, "", http.StatusNotFound)
			return nil
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
			if fromDateStr != "" && fromDate > (*visit.VisitedAt) {
				continue
			}
			if toDateStr != "" && toDate < (*visit.VisitedAt) {
				continue
			}
			if country != "" && country != (*location.Country) {
				continue
			}
			if toDistStr != "" && uint32(toDistance) <= (*location.Distance) {
				continue
			}
			var item schema.ResponceUserVisit
			item.Mark = visit.Mark
			item.Visited_at = visit.VisitedAt
			item.Place = location.Place
			resp.Visits = append(resp.Visits, &item)
		}

		sort.Sort(&resp)
		b, err := json.Marshal(resp)
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Content-Length", strconv.Itoa(len(b)))
		w.Write(b)

		return nil
	}, workers.TimeOut)

	checkTimeout(w, err)
}

func UpdateUserEndpoint(w http.ResponseWriter, req *http.Request) {
	_, err := workers.Wp.AddTaskSyncTimed(func() interface{} {
		id, err := parseID(req)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return nil
		}

		bytes, err := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return nil
		}

		ok := CheckNull(bytes)
		if ok {
			http.Error(w, "", http.StatusBadRequest)
			return nil
		}
		var user model.User
		_ = json.Unmarshal(bytes, &user)

		err = storage.DataStorage.User.Update(id, &user)
		if err != nil {

			if err == storage.ErrDoesNotExist {
				http.Error(w, "", http.StatusNotFound)
			} else {
				http.Error(w, "", http.StatusBadRequest)
			}
			return nil
		}

		w.Write([]byte("{}"))
		return nil
	}, workers.TimeOut)

	checkTimeout(w, err)
}

func CreateUserEndpoint(w http.ResponseWriter, req *http.Request) {
	_, err := workers.Wp.AddTaskSyncTimed(func() interface{} {
		var user model.User
		_ = json.NewDecoder(req.Body).Decode(&user)
		defer req.Body.Close()

		err := storage.DataStorage.User.Add(&user)
		if err != nil {

			if err == storage.ErrDoesNotExist {
				http.Error(w, "", http.StatusNotFound)
			} else {
				http.Error(w, "", http.StatusBadRequest)
			}
			return nil
		}

		w.Write([]byte("{}"))
		return nil
	}, workers.TimeOut)

	checkTimeout(w, err)
}
