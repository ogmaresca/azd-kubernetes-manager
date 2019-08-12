package templating

import (
	"text/template"

	"github.com/Masterminds/sprig"
)

// New creates a Go template with Sprig loaded
func New(name string) *template.Template {
	return template.New(name).Funcs(sprig.TxtFuncMap())
}
