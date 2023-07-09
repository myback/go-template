package template_functions

import (
	"fmt"
	"io"
	"path/filepath"
	"text/template"
)

func Render(file string, out io.Writer, values map[string]any) error {
	name := filepath.Base(file)

	tmpl, err := template.
		New(name).
		Funcs(createFuncMap()).
		Option("missingkey=error").
		ParseFiles(file)
	if err != nil {
		return fmt.Errorf("parse template failed: %s", err)
	}

	err = tmpl.Execute(out, values)
	if err != nil {
		return fmt.Errorf("render template failed: %s", err)
	}

	return nil
}
