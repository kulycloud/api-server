package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kulycloud/api-server/mapping"
	"net/http"
)

func (srv *Server) registerRouteRoutes(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/route").HandlerFunc(srv.getRoutes)
	router.Methods(http.MethodGet).Path("/route/{name}").HandlerFunc(srv.getRoute)
	router.Methods(http.MethodPost).Path("/route/{name}").HandlerFunc(srv.routeCreateUpdate(true))
	router.Methods(http.MethodPut).Path("/route/{name}").HandlerFunc(srv.routeCreateUpdate(false))
	router.Methods(http.MethodDelete).Path("/route/{name}").HandlerFunc(srv.deleteRoute)
}

type RouteListElement struct {
	Namespace	  string					`json:"namespace"`
	Name          string					`json:"name"`
	Specification *mapping.IncomingRoute	`json:"specification"`
}

func (srv *Server) getRoutes(w http.ResponseWriter, r *http.Request) {
	ctx := srv.getRequestContext()
	logger.Infow("Getting routes in namespace", "namespace", getNamespaceFromRequest(r))

	namespace := getNamespaceFromRequest(r)

	routeUids, err := srv.storage.GetRoutesInNamespace(ctx, namespace)
	if err != nil {
		logger.Error(err)
		writeError(w, http.StatusNotFound, ErrResourceNotFound)
		return
	}

	routes := make([]RouteListElement, 0, len(routeUids))
	for _, uid := range routeUids {
		r, err := srv.storage.GetRouteByUID(ctx, uid)
		if err != nil {
			writeError(w, http.StatusNotFound, ErrResourceNotFound)
			return
		}

		conv, err := mapping.MapRoute(r.Route)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		routes = append(routes, RouteListElement{
			Namespace:     r.Name.Namespace,
			Name:          r.Name.Name,
			Specification: conv,
		})
	}

	marshalled, err := json.Marshal(routes)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(marshalled)
}

func (srv *Server) getRoute(w http.ResponseWriter, r *http.Request) {
	ctx := srv.getRequestContext()

	logger.Infow("Getting route", "name", getNameFromRequest(r), "namespace", getNamespaceFromRequest(r))
	route, err := srv.storage.GetRouteByNamespacedName(ctx, getNamespaceFromRequest(r), getNameFromRequest(r))
	if err != nil {
		writeError(w, http.StatusNotFound, ErrResourceNotFound)
		return
	}

	conv, err := mapping.MapRoute(route.Route)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	marshalled, err := json.Marshal(conv)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(marshalled)
}

func (srv *Server) routeCreateUpdate(isCreate bool)  func(w http.ResponseWriter, r * http.Request) {

	return func(w http.ResponseWriter, r * http.Request) {
		ctx := srv.getRequestContext()

		namespace := getNamespaceFromRequest(r)
		name := getNameFromRequest(r)

		if !isNameValid(name) {
			writeError(w, http.StatusBadRequest, ErrInvalidName)
			return
		}

		_, err := srv.storage.GetRouteByNamespacedName(ctx, namespace, name)

		if isCreate {
			// POST Handler => may not exist
			if err == nil {
				writeError(w, http.StatusConflict, ErrResourceExists)
				return
			}
			logger.Infow("Creating route", "name", getNameFromRequest(r), "namespace", getNamespaceFromRequest(r))
		} else {
			// PUT Handler => has to exist
			if err != nil {
				writeError(w, http.StatusNotFound, ErrResourceNotFound)
				return
			}
			logger.Infow("Updating route", "name", getNameFromRequest(r), "namespace", getNamespaceFromRequest(r))
		}

		incomingRoute := &mapping.IncomingRoute{}

		err = json.NewDecoder(r.Body).Decode(incomingRoute)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		route, err := mapping.MapIncomingRoute(incomingRoute)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		_, err = srv.storage.SetRouteByNamespacedName(ctx, namespace, name, route)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		err = dispatchEvent(ResourceTypeRoute, getNamespaceFromRequest(r), getNameFromRequest(r))
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		conv, err := mapping.MapRoute(route)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		marshalled, err := json.Marshal(conv)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(marshalled)
	}
}

func (srv *Server) deleteRoute(w http.ResponseWriter, r *http.Request) {
	ctx := srv.getRequestContext()

	namespace := getNamespaceFromRequest(r)
	name := getNameFromRequest(r)

	route, err := srv.storage.GetRouteByNamespacedName(ctx, namespace, name)
	if err != nil {
		writeError(w, http.StatusNotFound, ErrResourceNotFound)
		return
	}

	logger.Infow("Deleting route", "name", getNameFromRequest(r), "namespace", getNamespaceFromRequest(r))

	err = srv.storage.DeleteRoute(ctx, namespace, name)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	err = dispatchEvent(ResourceTypeRoute, getNamespaceFromRequest(r), getNameFromRequest(r))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	conv, err := mapping.MapRoute(route.Route)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	marshalled, err := json.Marshal(conv)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(marshalled)
}
