package handler

import (
	"encoding/json"
	"net/http"

	"github.com/elBroom/highloadCup/app/model"
	"github.com/elBroom/highloadCup/app/schema"
	"github.com/elBroom/highloadCup/app/storage"
	"github.com/golang/glog"
)

func GetUserEndpoint(w http.ResponseWriter, req *http.Request) {
	glog.Infoln(req.Method, req.RequestURI)

	id, err := parseID(req)
	if err != nil {
		glog.Infoln(err)
		http.Error(w, "", http.StatusNotFound)
		return
	}

	user, ok := storage.DataStorage.User.Get(id)
	if !ok {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func VisitUserEndpoint(w http.ResponseWriter, req *http.Request) {
	glog.Infoln(req.Method, req.RequestURI)

	id, err := parseID(req)
	if err != nil {
		glog.Infoln(err)
		http.Error(w, "", http.StatusNotFound)
		return
	}

	var data schema.RequestUserVisits
	_ = json.NewDecoder(req.Body).Decode(&data)
	defer req.Body.Close()

	visits, ok := storage.DataStorage.VisitList.GetByUser(id)
	if !ok {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	var resp schema.ResponceUserVisits
	for _, visit := range visits {
		// TODO filter visit
		if data.Country != nil && (*data.Country) == (*visit.Location.Country) {

		}
		var item schema.ResponceUserVisit
		item.Mark = visit.Mark
		item.Visited_at = visit.VisitedAt
		item.Place = visit.Location.Place
		resp.Visits = append(resp.Visits, &item)
	}

	json.NewEncoder(w).Encode(resp)
}

func UpdateUserEndpoint(w http.ResponseWriter, req *http.Request) {
	glog.Infoln(req.Method, req.RequestURI)

	id, err := parseID(req)
	if err != nil {
		glog.Infoln(err)
		http.Error(w, "", http.StatusNotFound)
		return
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
		return
	}

	w.Write([]byte("{}"))
}

func CreateUserEndpoint(w http.ResponseWriter, req *http.Request) {
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
		return
	}

	w.Write([]byte("{}"))
}
