package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	protoStorage "github.com/kulycloud/protocol/storage"
	"net/http"
)

func (srv *Server) registerServiceRoutes(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/service").HandlerFunc(srv.getServices)
	router.Methods(http.MethodGet).Path("/service/{Name}").HandlerFunc(srv.getService)
}

type ServiceListElement struct {
	Namespace	  string				`json:"namespace"`
	Name          string				`json:"name"`
	Specification *protoStorage.Service	`json:"specification"`
}

func (srv *Server) getServices(w http.ResponseWriter, r *http.Request) {
	ctx := srv.getRequestContext()
	logger.Infow("Getting services in namespace", "namespace", getNamespaceFromRequest(r))

	namespace := getNamespaceFromRequest(r)

	serviceNames, err := srv.storage.GetServicesInNamespace(ctx, namespace)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	services := make([]ServiceListElement, 0, len(serviceNames))
	for _, name := range serviceNames {
		svc, err := srv.storage.GetService(ctx, namespace, name)
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}

		services = append(services, ServiceListElement{
			Namespace:     namespace,
			Name:          name,
			Specification: svc,
		})
	}

	marshalled, err := json.Marshal(services)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/marshalled")
	_, _ = w.Write(marshalled)

}

func (srv *Server) getService(w http.ResponseWriter, r *http.Request) {
	ctx := srv.getRequestContext()

	logger.Infow("Getting service", "Name", getNameFromRequest(r), "namespace", getNamespaceFromRequest(r))
	route, err := srv.storage.GetService(ctx, getNamespaceFromRequest(r), getNameFromRequest(r))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	marshalled, err := json.Marshal(route)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/marshalled")
	_, _ = w.Write(marshalled)
}

