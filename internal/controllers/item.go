package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"shop-api/internal/model"
	"shop-api/internal/storage"
	"shop-api/internal/views"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	//DefaultRenderFormat ...
	DefaultRenderFormat = "json"
	//ContentTypeHTML ...
	ContentTypeHTML = "text/html"
	//ContentTypeJSON ...
	ContentTypeJSON = "application/json"
)

//ItemController ...
type ItemController struct {
	logger  *logrus.Logger
	storage storage.Storage
}

//NewItemController ...
func NewItemController(logger *logrus.Logger, storage storage.Storage) *ItemController {
	return &ItemController{
		logger:  logger,
		storage: storage,
	}
}

//GetAllItems ...
func (ih *ItemController) GetAllItems(w http.ResponseWriter, r *http.Request) {
	ih.logger.Debugln("call GetAllItem handler")

	items := ih.storage.Items().GetItems()

	renderer, err := ih.getRendereByRequest(r)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	content, err := DoRender(renderer, items)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	if renderer.GetRendererType() == views.HTMLRendererType {
		w.Header().Add("Content-Type", ContentTypeHTML)
	} else {
		w.Header().Add("Content-Type", ContentTypeJSON)
	}
	respond(w, http.StatusOK, content)
}

//GetItem ...
func (ih *ItemController) GetItem(w http.ResponseWriter, r *http.Request) {
	ih.logger.Debugln("call GetItem handler")

	id := mux.Vars(r)["id"]
	intID, _ := strconv.Atoi(id)

	item := ih.storage.Items().GetItem(model.ItemID(intID))
	if item == nil {
		respondError(w, http.StatusNotFound, ErrItemNotFound)
		return
	}

	renderer, err := ih.getRendereByRequest(r)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	content, err := DoRender(renderer, item)
	if err != nil {
		respondError(w, http.StatusInternalServerError, errors.New("Internal server error"))
		return
	}

	if renderer.GetRendererType() == views.HTMLRendererType {
		w.Header().Add("Content-Type", ContentTypeHTML)
	} else {
		w.Header().Add("Content-Type", ContentTypeJSON)
	}
	respond(w, http.StatusOK, content)
}

//PutItem ...
func (ih *ItemController) PutItem(w http.ResponseWriter, r *http.Request) {
	ih.logger.Debugln("call PutItem handler")

	if r.Body == nil {
		respondError(w, http.StatusBadRequest, errors.New("empty request body"))
		return
	}

	item := &model.Item{}
	if err := json.NewDecoder(r.Body).Decode(item); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	if item.IsEmpty() {
		respondError(w, http.StatusBadRequest, ErrEmptyItem)
		return
	}

	item, err := ih.storage.Items().PutItem(item)
	if err != nil {
		respondError(w, http.StatusInternalServerError, errors.New("Internal server error"))
		return
	}

	renderer, err := ih.getRendereByRequest(r)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	content, err := DoRender(renderer, item)
	if err != nil {
		respondError(w, http.StatusInternalServerError, errors.New("Internal server error"))
		return
	}

	if renderer.GetRendererType() == views.HTMLRendererType {
		w.Header().Add("Content-Type", ContentTypeHTML)
	} else {
		w.Header().Add("Content-Type", ContentTypeJSON)
	}
	respond(w, http.StatusCreated, content)
}

//DeleteItem ...
func (ih *ItemController) DeleteItem(w http.ResponseWriter, r *http.Request) {
	ih.logger.Debugln("call DeleteItem handler")

	id, _ := mux.Vars(r)["id"]
	intID, _ := strconv.Atoi(id)

	item := ih.storage.Items().GetItem(model.ItemID(intID))
	if item == nil {
		respondError(w, http.StatusNotFound, ErrItemNotFound)
		return
	}

	if err := ih.storage.Items().DeleteItem(model.ItemID(intID)); err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respond(w, http.StatusNoContent, nil)
}

func respondError(w http.ResponseWriter, code int, err error) {
	type errorRespond struct {
		Err string `json:"error,omitempty"`
	}

	er := errorRespond{
		Err: err.Error(),
	}

	rawJSON, _ := json.Marshal(er)
	w.Header().Add("Content-Type", "application/json")

	respond(w, code, rawJSON)
}

func respond(w http.ResponseWriter, code int, content []byte) {
	w.WriteHeader(code)
	fmt.Fprint(w, string(content))
}

//DoRender ...
func DoRender(renderer views.Renderer, data interface{}) ([]byte, error) {
	return renderer.Render(data)
}

func (ih *ItemController) getRendereByRequest(r *http.Request) (views.Renderer, error) {
	switch getRenderFormat(r) {
	case "html":
		return views.NewHTMLRenderer(ih.logger)
	case "json":
		return views.NewJSONRenderer(ih.logger)
	default:
		return views.NewJSONRenderer(ih.logger)
	}
}

func getRenderFormat(r *http.Request) string {
	format := r.URL.Query().Get("format")
	if format == "" {
		format = DefaultRenderFormat
	}

	return format
}
