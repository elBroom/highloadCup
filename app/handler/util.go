package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func NotFound(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "", http.StatusBadRequest)
}

func parseID(req *http.Request) (uint32, error) {
	errParse := errors.New("Could not parse Id from request")
	vars := mux.Vars(req)
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
