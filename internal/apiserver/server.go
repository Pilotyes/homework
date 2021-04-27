package apiserver

import (
	"fmt"
	"net/http"
	"shop-api/internal/config"
	"shop-api/internal/controllers"
	"shop-api/internal/storage"
	"shop-api/internal/storage/mapstorage"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//Server ...
type Server struct {
	Config         *config.Config
	ItemController *controllers.ItemController
	Logger         *logrus.Logger
	Router         *mux.Router
	Storage        storage.Storage
}

//New ...
func New(mainConfig *config.Config) (*Server, error) {
	logger := logrus.New()

	level, err := logrus.ParseLevel(mainConfig.Server.LogLevel)
	if err != nil {
		return nil, err
	}

	logger.SetLevel(level)

	logger.Debugf("Created new server with next config sections:\n")
	logger.Debugf("\tServer: %#v\n", mainConfig.Server)
	logger.Debugf("\tDatabases: %#v\n", mainConfig.Databases)
	logger.Debugf("\tMongoDB: %#v\n", mainConfig.MongoDB)
	logger.Debugf("\tInternal: %#v\n", mainConfig.Internal)

	databaseDriver := ""
	for _, driver := range config.DatabaseDrivers {
		if mainConfig.Databases.Driver == driver {
			databaseDriver = mainConfig.Databases.Driver
		}
	}

	if databaseDriver == "" {
		return nil, fmt.Errorf("%s database not supported yet", mainConfig.Databases.Driver)
	}

	var storage storage.Storage
	switch databaseDriver {
	case config.MongoDBDriver:
		return nil, fmt.Errorf("Driver %s not implemented yet", config.MongoDBDriver)
	case config.InternalDriver:
		storage = mapstorage.New(mainConfig)
	}

	itemController := controllers.NewItemController(logger, storage)

	server := &Server{
		Config:         mainConfig,
		ItemController: itemController,
		Logger:         logger,
		Router:         mux.NewRouter(),
		Storage:        storage,
	}

	server.configureRouter()

	return server, nil
}

func (s *Server) configureRouter() {
	s.Router.Use(s.logRequest)

	handlers := []struct {
		path   string
		fn     func(http.ResponseWriter, *http.Request)
		method string
	}{
		{
			path:   "/items",
			fn:     s.ItemController.GetAllItems,
			method: http.MethodGet,
		},
		{
			path:   "/items",
			fn:     s.ItemController.PutItem,
			method: http.MethodPost,
		},
		{
			path:   "/items/{id:[0-9]+}",
			fn:     s.ItemController.GetItem,
			method: http.MethodGet,
		},
		{
			path:   "/items/{id:[0-9]+}",
			fn:     s.ItemController.DeleteItem,
			method: http.MethodDelete,
		},
	}

	for _, handler := range handlers {
		s.Router.HandleFunc(handler.path, handler.fn).Methods(handler.method)
		logger := s.Logger.WithFields(logrus.Fields{
			"path":   handler.path,
			"method": handler.method,
		})
		logger.Debugln("Registered new handler")
	}
}

func (s *Server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.Logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
		})
		logger.Infof("Recieved request %s %s", r.RequestURI, r.Method)

		start := time.Now()
		newResponseWriter := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(newResponseWriter, r)

		logger.Infof("result code = %d in = %f sec", newResponseWriter.code, time.Now().Sub(start).Seconds())
	})
}

//Start ...
func (s *Server) Start() error {
	s.Logger.Infof("Started listening server on address \"%s\"", s.Config.Server.BindAddr)
	return http.ListenAndServe(s.Config.Server.BindAddr, s.Router)
}
