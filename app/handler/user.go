package handler

import (
	"encoding/json"
	"net/http"

	"sort"

	"github.com/elBroom/highloadCup/app/model"
	"github.com/elBroom/highloadCup/app/schema"
	"github.com/elBroom/highloadCup/app/storage"
	"github.com/elBroom/highloadCup/app/workers"
	"github.com/golang/glog"
)

func GetUserEndpoint(w http.ResponseWriter, req *http.Request) {
	_, err := workers.Wp.AddTaskSyncTimed(func() interface{} {
		glog.Infoln(req.Method, req.RequestURI)

		id, err := parseID(req)
		if err != nil {
			glog.Infoln(err)
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
		glog.Infoln(req.Method, req.RequestURI)

		id, err := parseID(req)
		if err != nil {
			glog.Infoln(err)
			http.Error(w, "", http.StatusNotFound)
			return nil
		}

		var data schema.RequestUserVisits
		_ = json.NewDecoder(req.Body).Decode(&data)
		defer req.Body.Close()

		visits, ok := storage.DataStorage.VisitList.GetByUser(id)
		if !ok {
			http.Error(w, "", http.StatusNotFound)
			return nil
		}

		var resp schema.ResponceUserVisits
		for _, visit := range visits {
			if data.FromDate != nil && (*data.FromDate) >= (*visit.VisitedAt) {
				continue
			}
			if data.ToDate != nil && (*data.ToDate) <= (*visit.VisitedAt) {
				continue
			}
			if data.Country != nil && (*data.Country) != (*visit.Location.Country) {
				continue
			}
			if data.ToDistance != nil && (*data.ToDistance) >= (*visit.Location.Distance) {
				continue
			}
			var item schema.ResponceUserVisit
			item.Mark = visit.Mark
			item.Visited_at = visit.VisitedAt
			item.Place = visit.Location.Place
			resp.Visits = append(resp.Visits, &item)
		}

		sort.Sort(&resp)
		json.NewEncoder(w).Encode(resp)
		return nil
	}, workers.TimeOut)

	checkTimeout(w, err)
}

func UpdateUserEndpoint(w http.ResponseWriter, req *http.Request) {
	_, err := workers.Wp.AddTaskSyncTimed(func() interface{} {
		glog.Infoln(req.Method, req.RequestURI)

		id, err := parseID(req)
		if err != nil {
			glog.Infoln(err)
			http.Error(w, "", http.StatusNotFound)
			return nil
		}

		var user model.User
		_ = json.NewDecoder(req.Body).Decode(&user)
		defer req.Body.Close()

		err = storage.DataStorage.User.Update(id, &user)
		if err != nil {
			glog.Infoln(err)
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
		glog.Infoln(req.Method, req.RequestURI)

		var user model.User
		_ = json.NewDecoder(req.Body).Decode(&user)
		defer req.Body.Close()

		err := storage.DataStorage.User.Add(&user)
		if err != nil {
			glog.Infoln(err)
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
