package views

const (
	//JSONRendererType ...
	JSONRendererType = "json"
	//HTMLRendererType ...
	HTMLRendererType = "html"
)

//Renderer ...
type Renderer interface {
	Render(data interface{}) ([]byte, error)
	GetRendererType() string
}
