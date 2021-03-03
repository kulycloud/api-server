package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func getNamespaceFromRequest(r *http.Request) string {
	return mux.Vars(r)["namespace"]
}

func getNameFromRequest(r *http.Request) string {
	return mux.Vars(r)["Name"]
}

const MarshallErrorString = `{"error": "could not marshall error description"}`

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
