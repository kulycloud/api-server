package server

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

var ErrResourceNotFound = errors.New("resource not found")
var ErrResourceExists = errors.New("resource already exists")
var ErrInvalidName = errors.New("invalid name")

const MarshallErrorString = `{"error": "could not marshall error description"}`

func getNamespaceFromRequest(r *http.Request) string {
	return mux.Vars(r)["namespace"]
}

func getNameFromRequest(r *http.Request) string {
	return mux.Vars(r)["name"]
}

type ErrorResponse struct {
	Text string `json:"error"`
}

func writeError(w http.ResponseWriter, code int, err error) {
	resp := ErrorResponse{
		Text: err.Error(),
	}
	w.Header().Set("Content-Type", "application/json")

	marshalled, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(MarshallErrorString))
		return
	}

	w.WriteHeader(code)
	_, _ = w.Write(marshalled)
}

func isNameValid(name string) bool {
	return len(name) != 0
}
