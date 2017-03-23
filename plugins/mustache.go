package main

import (
	"fmt"

	"github.com/cbroglie/mustache"
)

// Execute the template
func Execute(input string, data interface{}) ([]byte, error) {
	var t *mustache.Template
	var err error

	if t, err = mustache.ParseFile(input); err != nil {
		return nil, fmt.Errorf("Could not parse template %s", input)
	}

	o, err := t.Render(data)
	if err != nil {
		return nil, err
	}

	return []byte(o), nil
}
