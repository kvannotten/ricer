package main

import (
	"bytes"
	"fmt"
	"html/template"
)

// Execute the template
func Execute(input string, data interface{}) ([]byte, error) {
	var t *template.Template
	var err error

	if t, err = template.ParseFiles(input); err != nil {
		return nil, fmt.Errorf("Could not parse template %s", input)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
