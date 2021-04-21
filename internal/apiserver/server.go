package apiserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"shop-api/internal/config"
	"shop-api/internal/model"
	"shop-api/internal/storage"
	"shop-api/internal/storage/mapstorage"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//Server ...
type Server struct {
	config  *config.Config
	logger  *logrus.Logger
	router  *mux.Router
	storage storage.Storage
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

	server := &Server{
		config:  mainConfig,
		logger:  logger,
		router:  mux.NewRouter(),
		storage: storage,
	}

	server.configureRouter()

	return server, nil
}

func (s *Server) configureRouter() {
	s.router.Use(s.logRequest)

	handlers := map[string]struct {
		fn      func(http.ResponseWriter, *http.Request)
		methods []string
	}{
		"/items": {
			s.handlerItems,
			[]string{
				http.MethodGet,
				http.MethodPost,
			},
		},
		"/items/{id:[0-9]+}": {
			s.handlerItems,
			[]string{
				http.MethodGet,
				http.MethodDelete,
			},
		},
	}

	for path, handler := range handlers {
		s.router.HandleFunc(path, handler.fn).Methods(handler.methods...)
		logger := s.logger.WithFields(logrus.Fields{
			"path":    path,
			"methods": handler.methods,
		})
		logger.Debugln("Registered new handler")
	}
}

func (s *Server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
		})
		logger.Infof("Recieved request %s %s", r.RequestURI, r.Method)

		start := time.Now()
		newResponseWriter := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(newResponseWriter, r)

		logger.Infof("result code = %d in = %f sec", newResponseWriter.code, time.Now().Sub(start).Seconds())
	})
}

func (s *Server) handlerItems(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		vars := mux.Vars(r)
		if id, ok := vars["id"]; ok {
			intID, err := strconv.Atoi(id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, err)
				return
			}

			item := s.storage.Items().GetItem(model.ItemID(intID))
			if item == nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintln(w, errors.New("Item not found"))
				return
			}

			rawJSON, err := json.Marshal(item)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, err)
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, string(rawJSON))
		} else {
			items := s.storage.Items().GetItems()

			rawJSON, err := json.Marshal(items)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, err)
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, string(rawJSON))
		}
	case http.MethodPost:
		if r.Body == nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, errors.New("Empty body"))
			return
		}

		item := &model.Item{}
		if err := json.NewDecoder(r.Body).Decode(item); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, err)
			return
		}

		if item.IsEmpty() {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, errors.New("Item must not be empty"))
			return
		}

		item, err := s.storage.Items().PutItem(item)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(item)
	case http.MethodDelete:
		vars := mux.Vars(r)
		if id, ok := vars["id"]; ok {
			intID, err := strconv.Atoi(id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, err)
				return
			}

			if err := s.storage.Items().DeleteItem(model.ItemID(intID)); err != nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintln(w, err)
				return
			}

			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}
}

//Start ...
func (s *Server) Start() error {
	s.logger.Infof("Started listening server on address \"%s\"", s.config.Server.BindAddr)
	return http.ListenAndServe(s.config.Server.BindAddr, s.router)
}
