package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kulycloud/api-server/communication"
	commonCommunication "github.com/kulycloud/common/communication"
	"net/http"
)

var ErrResourceNotFound = errors.New("resource not found")
var ErrResourceExists = errors.New("resource already exists")
var ErrInvalidName = errors.New("invalid name")

const MarshallErrorString = `{"error": "could not marshall error description"}`

const ResourceTypeService = "service"
const ResourceTypeRoute = "route"

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

func dispatchEvent(resourceType, namespace, name string) error {
	evt := commonCommunication.NewConfigurationChanged(commonCommunication.NewResource(
		resourceType,
		name,
		namespace,
	))

	err := communication.ControlPlane.CreateEvent(evt)
	if err != nil {
		return fmt.Errorf("could not dispatch event: %w", err)
	}
	return nil
}
