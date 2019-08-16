package templating

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig"

	"github.com/alexcesaro/log/stdlog"
)

var (
	logger = stdlog.GetFromFlags()
)

// New creates a Go template with Sprig loaded
func New(name string) *template.Template {
	return template.New(name).Funcs(sprig.TxtFuncMap())
}

// Execute a template with the given args
func Execute(name string, value string, args interface{}) (string, error) {
	if value == "" {
		return value, nil
	}
	template := New(name)
	template, err := template.Parse(value)
	if err != nil {
		return value, fmt.Errorf("Error parsing template: %s", err.Error())
	}
	buffer := new(bytes.Buffer)
	err = template.Execute(buffer, args)
	if err != nil {
		return value, fmt.Errorf("Error executing template: %s", err.Error())
	}
	return buffer.String(), nil
}
