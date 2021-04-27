package views

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

//JSONRenderer ...
type JSONRenderer struct {
	logger       *logrus.Entry
	rendererType string
}

//NewJSONRenderer ...
func NewJSONRenderer(logger *logrus.Logger) (*JSONRenderer, error) {
	return &JSONRenderer{
		logger: logger.WithFields(logrus.Fields{
			"renderer": "json",
		}),
		rendererType: JSONRendererType,
	}, nil
}

//Render ...
func (jr *JSONRenderer) Render(data interface{}) ([]byte, error) {
	jr.logger.Debugln("Rendering data", data)

	rawJSON, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json renderer error: %w", err)
	}

	return rawJSON, nil
}

//GetRendererType ...
func (jr *JSONRenderer) GetRendererType() string {
	return jr.rendererType
}
