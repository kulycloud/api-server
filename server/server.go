package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kulycloud/api-server/config"
	"github.com/kulycloud/common/communication"
	"github.com/kulycloud/common/logging"
	"net/http"
)

var logger = logging.GetForComponent("server")

type Server struct {
	router *mux.Router
	storage *communication.StorageCommunicator
}

func NewServer(storage *communication.StorageCommunicator) *Server {
	srv := &Server{
		router:  mux.NewRouter(),
		storage: storage,
	}

	srv.router.Use(srv.storageAvailabilityMiddleware)

	namespacedRouter := srv.router.PathPrefix("/{namespace}/").Subrouter()

	srv.registerServiceRoutes(namespacedRouter)
	srv.registerRouteRoutes(namespacedRouter)

	return srv
}

func (srv *Server) Start() error {
	logger.Info("Starting HTTP listener")
	return http.ListenAndServe(fmt.Sprintf(":%v", config.GlobalConfig.HTTPPort), srv.router)
}

func (srv *Server) storageAvailabilityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !srv.storage.Ready() {
			w.WriteHeader(500)
			w.Write([]byte("storage is not ready"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (srv *Server) getRequestContext() context.Context {
	return context.Background()
}
