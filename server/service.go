package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	protoStorage "github.com/kulycloud/protocol/storage"
	"net/http"
)

func (srv *Server) registerServiceRoutes(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/service").HandlerFunc(srv.getServices)
	router.Methods(http.MethodGet).Path("/service/{name}").HandlerFunc(srv.getService)
	router.Methods(http.MethodPost).Path("/service/{name}").HandlerFunc(srv.serviceCreateUpdate(true))
	router.Methods(http.MethodPut).Path("/service/{name}").HandlerFunc(srv.serviceCreateUpdate(false))
	router.Methods(http.MethodDelete).Path("/service/{name}").HandlerFunc(srv.deleteService)
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
		writeError(w, http.StatusNotFound, ErrResourceNotFound)
		return
	}

	services := make([]ServiceListElement, 0, len(serviceNames))
	for _, name := range serviceNames {
		svc, err := srv.storage.GetService(ctx, namespace, name)
		if err != nil {
			writeError(w, http.StatusNotFound, ErrResourceNotFound)
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

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(marshalled)

}

func (srv *Server) getService(w http.ResponseWriter, r *http.Request) {
	ctx := srv.getRequestContext()

	logger.Infow("Getting service", "name", getNameFromRequest(r), "namespace", getNamespaceFromRequest(r))
	service, err := srv.storage.GetService(ctx, getNamespaceFromRequest(r), getNameFromRequest(r))
	if err != nil {
		writeError(w, http.StatusNotFound, ErrResourceNotFound)
		return
	}

	marshalled, err := json.Marshal(service)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(marshalled)
}

func (srv *Server) serviceCreateUpdate(isCreate bool)  func(w http.ResponseWriter, r * http.Request) {

	return func(w http.ResponseWriter, r * http.Request) {
		ctx := srv.getRequestContext()

		namespace := getNamespaceFromRequest(r)
		name := getNameFromRequest(r)

		if len(name) == 0 {
			writeError(w, http.StatusBadRequest, ErrInvalidName)
			return
		}

		_, err := srv.storage.GetService(ctx, namespace, name)

		if isCreate {
			// POST Handler => may not exist
			if err == nil {
				writeError(w, http.StatusConflict, ErrResourceExists)
				return
			}
			logger.Infow("Creating service", "name", getNameFromRequest(r), "namespace", getNamespaceFromRequest(r))
		} else {
			// PUT Handler => has to exist
			if err != nil {
				writeError(w, http.StatusNotFound, ErrResourceNotFound)
				return
			}
			logger.Infow("Updating service", "name", getNameFromRequest(r), "namespace", getNamespaceFromRequest(r))
		}

		service := &protoStorage.Service{}

		err = json.NewDecoder(r.Body).Decode(service)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		err = srv.storage.SetService(ctx, namespace, name, service)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		marshalled, err := json.Marshal(service)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(marshalled)
	}
}

func (srv *Server) deleteService(w http.ResponseWriter, r *http.Request) {
	ctx := srv.getRequestContext()

	namespace := getNamespaceFromRequest(r)
	name := getNameFromRequest(r)

	service, err := srv.storage.GetService(ctx, namespace, name)
	if err != nil {
		writeError(w, http.StatusNotFound, ErrResourceNotFound)
		return
	}

	logger.Infow("Deleting service", "name", getNameFromRequest(r), "namespace", getNamespaceFromRequest(r))

	err = srv.storage.DeleteService(ctx, namespace, name)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	marshalled, err := json.Marshal(service)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(marshalled)
}
