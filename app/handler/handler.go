package handler

import (
	"errors"
	"net/http"
	"strconv"

	"encoding/json"

	"github.com/elBroom/highloadCup/app/schema"
	"github.com/elBroom/highloadCup/app/storage"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

func GetUserEndpoint(w http.ResponseWriter, req *http.Request) {
	glog.Infoln(req.Method, req.RequestURI)

	id, err := parseID(req)
	if err != nil {
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

	for _, visit := range visits {
		// TODO filter visit
		if data.Country != nil && (*data.Country) == (*visit.Location.Country) {

		}
	}

	var resp schema.ResponceUserVisits
	json.NewEncoder(w).Encode(resp)
}

func GetLocatioAvgnEndpoint(w http.ResponseWriter, req *http.Request) {
	glog.Infoln(req.Method, req.RequestURI)

	id, err := parseID(req)
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	var data schema.RequestLocationVisits
	_ = json.NewDecoder(req.Body).Decode(&data)
	defer req.Body.Close()

	visits, ok := storage.DataStorage.VisitList.GetByLocation(id)
	if !ok {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	for _, visit := range visits {
		// TODO filter visit
		if data.Gender != nil && (*data.Gender) == (*visit.User.Gender) {

		}
	}

	var resp schema.ResponceLocationVisits
	json.NewEncoder(w).Encode(resp)
}

func GetLocationEndpoint(w http.ResponseWriter, req *http.Request) {
	glog.Infoln(req.Method, req.RequestURI)

	id, err := parseID(req)
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	user, ok := storage.DataStorage.Location.Get(id)
	if !ok {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func GetVisitEndpoint(w http.ResponseWriter, req *http.Request) {
	glog.Infoln(req.Method, req.RequestURI)

	id, err := parseID(req)
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	user, ok := storage.DataStorage.Location.Get(id)
	if !ok {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func parseID(r *http.Request) (uint32, error) {
	errParse := errors.New("Could not parse Id from request")
	vars := mux.Vars(r)
	strID, ok := vars["id"]
	if !ok {
		return 0, errParse
	}
	id64, err := strconv.ParseInt(strID, 10, 32)
	id := uint32(id64)
	if err != nil {
		return 0, errParse
	}
	return id, nil
}
