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
			ok := storage.DataStorage.Visit.FetchLocation(visit, storage.DataStorage)
			if !ok {
				continue
			}
			if fromDateStr != "" && visit.VisitedAt != nil && fromDate > (*visit.VisitedAt) {
				continue
			}
			if toDateStr != "" && visit.VisitedAt != nil && toDate < (*visit.VisitedAt) {
				continue
			}
			if country != "" && visit.Location.Country != nil && country != (*visit.Location.Country) {
				continue
			}
			if toDistStr != "" && visit.Location.Distance != nil && uint32(toDistance) <= (*visit.Location.Distance) {
				continue
			}
			if visit.VisitedAt != nil {
				var item schema.ResponceUserVisit
				item.Mark = visit.Mark
				item.Visited_at = visit.VisitedAt
				item.Place = visit.Location.Place
				resp.Visits = append(resp.Visits, &item)
			}
		}

		sort.Sort(&resp)
		json.NewEncoder(w).Encode(resp)
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
