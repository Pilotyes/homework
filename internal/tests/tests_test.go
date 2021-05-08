package test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"shop-api/internal/apiserver"
	"shop-api/internal/config"
	"shop-api/internal/controllers"
	"testing"

	"github.com/stretchr/testify/assert"
)

type test struct {
	path            string
	method          string
	requestBody     io.Reader
	responseBody    string
	responseHeaders http.Header
	statusCode      int
}

//IntegrateTests ...
func Test_integrate_tests(t *testing.T) {
	mainConfig := config.NewConfig()
	mainConfig.Server.LogLevel = "debug"

	server, err := apiserver.New(mainConfig)
	if err != nil {
		t.Fatal(err)
	}

	tests := []test{
		{
			path:         "/items",
			method:       http.MethodGet,
			requestBody:  nil,
			responseBody: `[]`,
			responseHeaders: http.Header{
				"Content-Type": {controllers.ContentTypeJSON},
			},
			statusCode: http.StatusOK,
		},
		{
			path:         "/items?format=html",
			method:       http.MethodGet,
			requestBody:  nil,
			responseBody: "some HTML",
			responseHeaders: http.Header{
				"Content-Type": {controllers.ContentTypeHTML},
			},
			statusCode: http.StatusOK,
		},
		{
			path:         "/items",
			method:       http.MethodPost,
			requestBody:  bytes.NewBuffer([]byte(`{}`)),
			responseBody: marshalError(controllers.ErrEmptyItem),
			responseHeaders: http.Header{
				"Content-Type": {controllers.ContentTypeJSON},
			},
			statusCode: http.StatusBadRequest,
		},
		{
			path:         "/items",
			method:       http.MethodPost,
			requestBody:  nil,
			responseBody: marshalError(controllers.ErrEmptyRequestBody),
			responseHeaders: http.Header{
				"Content-Type": {controllers.ContentTypeJSON},
			},
			statusCode: http.StatusBadRequest,
		},
		{
			path:         "/items",
			method:       http.MethodGet,
			requestBody:  nil,
			responseBody: `[]`,
			responseHeaders: http.Header{
				"Content-Type": {controllers.ContentTypeJSON},
			},
			statusCode: http.StatusOK,
		},
		{
			path:         "/items/1",
			method:       http.MethodGet,
			requestBody:  nil,
			responseBody: marshalError(controllers.ErrItemNotFound),
			responseHeaders: http.Header{
				"Content-Type": {controllers.ContentTypeJSON},
			},
			statusCode: http.StatusNotFound,
		},
		{
			path:         "/items/0",
			method:       http.MethodGet,
			requestBody:  nil,
			responseBody: marshalError(controllers.ErrItemNotFound),
			responseHeaders: http.Header{
				"Content-Type": {controllers.ContentTypeJSON},
			},
			statusCode: http.StatusNotFound,
		},
		{
			path:         "/items/1",
			method:       http.MethodDelete,
			requestBody:  nil,
			responseBody: marshalError(controllers.ErrItemNotFound),
			responseHeaders: http.Header{
				"Content-Type": {controllers.ContentTypeJSON},
			},
			statusCode: http.StatusNotFound,
		},
		{
			path:         "/items",
			method:       http.MethodPost,
			requestBody:  bytes.NewBuffer([]byte(`{invalid json}`)),
			responseBody: marshalError(errors.New(`invalid character 'i' looking for beginning of object key string`)),
			responseHeaders: http.Header{
				"Content-Type": {controllers.ContentTypeJSON},
			},
			statusCode: http.StatusBadRequest,
		},
		{
			path:         "/items",
			method:       http.MethodPost,
			requestBody:  bytes.NewBuffer([]byte(`{"invalid_field_name": "invalid_fiels_value"}`)),
			responseBody: marshalError(controllers.ErrEmptyItem),
			responseHeaders: http.Header{
				"Content-Type": {controllers.ContentTypeJSON},
			},
			statusCode: http.StatusBadRequest,
		},
		{
			path:         "/items/1",
			method:       http.MethodDelete,
			requestBody:  nil,
			responseBody: marshalError(controllers.ErrItemNotFound),
			responseHeaders: http.Header{
				"Content-Type": {controllers.ContentTypeJSON},
			},
			statusCode: http.StatusNotFound,
		},
		{
			path:         "/items",
			method:       http.MethodPost,
			requestBody:  bytes.NewBuffer([]byte(`{"name":"name1","description":"description1","original_price":10,"discount_price":5,"articul":100001,"category":"category1"}`)),
			responseBody: `{"id":1,"name":"name1","description":"description1","original_price":10,"discount_price":5,"articul":100001,"category":"category1"}`,
			responseHeaders: http.Header{
				"Content-Type": {controllers.ContentTypeJSON},
			},
			statusCode: http.StatusCreated,
		},
		{
			path:         "/items",
			method:       http.MethodGet,
			requestBody:  nil,
			responseBody: `[{"id":1,"name":"name1","description":"description1","original_price":10,"discount_price":5,"articul":100001,"category":"category1"}]`,
			responseHeaders: http.Header{
				"Content-Type": {controllers.ContentTypeJSON},
			},
			statusCode: http.StatusOK,
		},
		{
			path:         "/items/2",
			method:       http.MethodGet,
			requestBody:  nil,
			responseBody: marshalError(controllers.ErrItemNotFound),
			responseHeaders: http.Header{
				"Content-Type": {controllers.ContentTypeJSON},
			},
			statusCode: http.StatusNotFound,
		},
		{
			path:         "/items/1",
			method:       http.MethodGet,
			requestBody:  nil,
			responseBody: `{"id":1,"name":"name1","description":"description1","original_price":10,"discount_price":5,"articul":100001,"category":"category1"}`,
			responseHeaders: http.Header{
				"Content-Type": {controllers.ContentTypeJSON},
			},
			statusCode: http.StatusOK,
		},
		{
			path:         "/items/001",
			method:       http.MethodGet,
			requestBody:  nil,
			responseBody: `{"id":1,"name":"name1","description":"description1","original_price":10,"discount_price":5,"articul":100001,"category":"category1"}`,
			responseHeaders: http.Header{
				"Content-Type": {controllers.ContentTypeJSON},
			},
			statusCode: http.StatusOK,
		},
		{
			path:         "/items/1?format=html",
			method:       http.MethodGet,
			requestBody:  nil,
			responseBody: "some HTML",
			responseHeaders: http.Header{
				"Content-Type": {controllers.ContentTypeHTML},
			},
			statusCode: http.StatusOK,
		},
		{
			path:         "/items/1?format=invalid",
			method:       http.MethodGet,
			requestBody:  nil,
			responseBody: `{"id":1,"name":"name1","description":"description1","original_price":10,"discount_price":5,"articul":100001,"category":"category1"}`,
			responseHeaders: http.Header{
				"Content-Type": {controllers.ContentTypeJSON},
			},
			statusCode: http.StatusOK,
		},
		{
			path:         "/items",
			method:       http.MethodPost,
			requestBody:  bytes.NewBuffer([]byte(`{"name":"name2","description":"description2","original_price":20,"discount_price":15,"articul":100002,"category":"category2"}`)),
			responseBody: `{"id":2,"name":"name2","description":"description2","original_price":20,"discount_price":15,"articul":100002,"category":"category2"}`,
			responseHeaders: http.Header{
				"Content-Type": {controllers.ContentTypeJSON},
			},
			statusCode: http.StatusCreated,
		},
		{
			path:         "/items?format=html",
			method:       http.MethodPost,
			requestBody:  bytes.NewBuffer([]byte(`{"description": "description3"`)),
			responseBody: marshalError(errors.New("unexpected EOF")),
			responseHeaders: http.Header{
				"Content-Type": {controllers.ContentTypeJSON},
			},
			statusCode: http.StatusBadRequest,
		},
		{
			path:         "/items?format=html",
			method:       http.MethodPost,
			requestBody:  bytes.NewBuffer([]byte(`{"name":"name3","description":"description3","original_price":50,"discount_price":40,"articul":100003,"category":"category3"}`)),
			responseBody: `some HTML`,
			responseHeaders: http.Header{
				"Content-Type": {controllers.ContentTypeHTML},
			},
			statusCode: http.StatusCreated,
		},
		{
			path:         "/items",
			method:       http.MethodPost,
			requestBody:  bytes.NewBuffer([]byte(`{"name":"name4","description":"description4","original_price":40,"discount_price":30,"articul":100004,"category":"category4"}`)),
			responseBody: `{"id":4,"name":"name4","description":"description4","original_price":40,"discount_price":30,"articul":100004,"category":"category4"}`,
			responseHeaders: http.Header{
				"Content-Type": {controllers.ContentTypeJSON},
			},
			statusCode: http.StatusCreated,
		},
		{
			path:         "/items",
			method:       http.MethodGet,
			requestBody:  nil,
			responseBody: `[{"id":1,"name":"name1","description":"description1","original_price":10,"discount_price":5,"articul":100001,"category":"category1"},{"id":2,"name":"name2","description":"description2","original_price":20,"discount_price":15,"articul":100002,"category":"category2"},{"id":3,"name":"name3","description":"description3","original_price":50,"discount_price":40,"articul":100003,"category":"category3"},{"id":4,"name":"name4","description":"description4","original_price":40,"discount_price":30,"articul":100004,"category":"category4"}]`,
			responseHeaders: http.Header{
				"Content-Type": {controllers.ContentTypeJSON},
			},
			statusCode: http.StatusOK,
		},
		{
			path:            "/items/1",
			method:          http.MethodDelete,
			requestBody:     nil,
			responseBody:    ``,
			responseHeaders: http.Header{},
			statusCode:      http.StatusNoContent,
		},
		{
			path:         "/items",
			method:       http.MethodGet,
			requestBody:  nil,
			responseBody: `[{"id":2,"name":"name2","description":"description2","original_price":20,"discount_price":15,"articul":100002,"category":"category2"},{"id":3,"name":"name3","description":"description3","original_price":50,"discount_price":40,"articul":100003,"category":"category3"},{"id":4,"name":"name4","description":"description4","original_price":40,"discount_price":30,"articul":100004,"category":"category4"}]`,
			responseHeaders: http.Header{
				"Content-Type": {controllers.ContentTypeJSON},
			},
			statusCode: http.StatusOK,
		},
		{
			path:         "/items/1",
			method:       http.MethodGet,
			requestBody:  nil,
			responseBody: marshalError(controllers.ErrItemNotFound),
			responseHeaders: http.Header{
				"Content-Type": {controllers.ContentTypeJSON},
			},
			statusCode: http.StatusNotFound,
		},
		{
			path:         "/items",
			method:       http.MethodPost,
			requestBody:  bytes.NewBuffer([]byte(`{"name":"name5","description":"description5","original_price":60,"discount_price":50,"articul":100005,"category":"category5"}`)),
			responseBody: `{"id":5,"name":"name5","description":"description5","original_price":60,"discount_price":50,"articul":100005,"category":"category5"}`,
			responseHeaders: http.Header{
				"Content-Type": {controllers.ContentTypeJSON},
			},
			statusCode: http.StatusCreated,
		},
		{
			path:            "/items/2",
			method:          http.MethodDelete,
			requestBody:     nil,
			responseBody:    ``,
			responseHeaders: http.Header{},
			statusCode:      http.StatusNoContent,
		},
		{
			path:            "/items/3",
			method:          http.MethodDelete,
			requestBody:     nil,
			responseBody:    ``,
			responseHeaders: http.Header{},
			statusCode:      http.StatusNoContent,
		},
		{
			path:            "/items/4",
			method:          http.MethodDelete,
			requestBody:     nil,
			responseBody:    ``,
			responseHeaders: http.Header{},
			statusCode:      http.StatusNoContent,
		},
		{
			path:         "/items",
			method:       http.MethodGet,
			requestBody:  nil,
			responseBody: `[{"id":5,"name":"name5","description":"description5","original_price":60,"discount_price":50,"articul":100005,"category":"category5"}]`,
			responseHeaders: http.Header{
				"Content-Type": {controllers.ContentTypeJSON},
			},
			statusCode: http.StatusOK,
		},
		{
			path:            "/items/5",
			method:          http.MethodDelete,
			requestBody:     nil,
			responseBody:    ``,
			responseHeaders: http.Header{},
			statusCode:      http.StatusNoContent,
		},
		{
			path:         "/items?format=html",
			method:       http.MethodGet,
			requestBody:  nil,
			responseBody: `some HTML`,
			responseHeaders: http.Header{
				"Content-Type": {controllers.ContentTypeHTML},
			},
			statusCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s_%s", tc.path, tc.method), func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.path, tc.requestBody)
			if err != nil {
				t.Fatal(err)
			}

			respRec := httptest.NewRecorder()
			server.Router.ServeHTTP(respRec, req)

			assert.Equal(t, tc.statusCode, respRec.Code)
			assert.Equal(t, tc.responseHeaders, respRec.Header())

			if respRec.Header().Get("Content-Type") != controllers.ContentTypeHTML {
				assert.Equal(t, tc.responseBody, respRec.Body.String())
			} else {
				assert.NotEqual(t, len(respRec.Body.String()), 0)
			}
		})
	}
}

func marshalError(err error) string {
	return fmt.Sprintf("{\"error\":\"%s\"}", err)
}
