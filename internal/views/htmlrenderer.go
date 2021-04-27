package views

import (
	"bytes"
	"fmt"
	"shop-api/internal/model"
	"text/template"

	"github.com/sirupsen/logrus"
)

var (
	templates = map[string]string{
		// "item_template":  "./internal/views/templates/item_template.html",
		// "items_template": "./internal/views/templates/items_template.html",
		"item_template":  "../../internal/views/templates/item_template.html",
		"items_template": "../../internal/views/templates/items_template.html",
	}
	globalTemplate *template.Template
)

//HTMLRenderer ...
type HTMLRenderer struct {
	logger       *logrus.Entry
	template     *template.Template
	rendererType string
}

//NewHTMLRenderer ...
func NewHTMLRenderer(logger *logrus.Logger) (*HTMLRenderer, error) {
	HTMLRenderer := &HTMLRenderer{
		logger: logger.WithFields(logrus.Fields{
			"renderer": "html",
		}),
		rendererType: HTMLRendererType,
	}

	if globalTemplate == nil {
		tmpl, err := template.ParseFiles(templates["item_template"], templates["items_template"])
		if err != nil {
			err := fmt.Errorf("parse templates error: %w", err)
			HTMLRenderer.logger.Errorln(err.Error())
			return nil, err
		}

		globalTemplate = tmpl
	}

	HTMLRenderer.template = globalTemplate

	return HTMLRenderer, nil
}

//Render ...
func (hr *HTMLRenderer) Render(data interface{}) ([]byte, error) {
	hr.logger.Debugln("Rendering data", data)

	buf := bytes.Buffer{}

	var err error
	if item, ok := data.(*model.Item); ok {
		err = hr.template.ExecuteTemplate(&buf, "item_template.html", map[string]interface{}{
			"item": item,
		})
	} else if items, ok := data.([]*model.Item); ok {
		err = hr.template.ExecuteTemplate(&buf, "items_template.html", map[string]interface{}{
			"items": items,
		})
	} else {
		err := fmt.Errorf("data not implemented items")
		hr.logger.Errorln(err.Error())
		return nil, err
	}

	if err != nil {
		err := fmt.Errorf("template execute error: %w", err)
		hr.logger.Errorln(err.Error())
		return nil, err
	}

	return buf.Bytes(), nil
}

//GetRendererType ...
func (hr *HTMLRenderer) GetRendererType() string {
	return hr.rendererType
}
