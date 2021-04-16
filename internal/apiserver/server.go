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
	config  *config.ServerConfig
	logger  *logrus.Logger
	router  *mux.Router
	storage storage.Storage
}

//New ...
func New(config *config.ServerConfig) (*Server, error) {
	logger := logrus.New()

	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		return nil, err
	}

	logger.SetLevel(level)

	logger.Debugf("Created new server with next config %#v\n", config)

	storage := mapstorage.New()

	server := &Server{
		config:  config,
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
		"/items/{id:[0-9]+}/": {
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

			item := s.storage.Items().GetItem(intID)
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

		item, err := s.storage.Items().PutItem(item)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(item)
	}
}

//Start ...
func (s *Server) Start() error {
	s.logger.Infof("Started listening server on address \"%s\"", s.config.BindAddr)
	return http.ListenAndServe(s.config.BindAddr, s.router)
}
