package main

import (
	"fmt"
	"os"

	"github.com/cbroglie/mustache"
)

// Execute the template
func Execute(input, output string, data interface{}) error {
	var t *mustache.Template
	var err error

	if t, err = mustache.ParseFile(input); err != nil {
		return fmt.Errorf("Could not parse template %s", input)
	}

	// write the output file
	var f *os.File
	if f, err = os.Create(output); err != nil {
		return fmt.Errorf("[2] Could not create %s from %s", output, input)
	}

	fmt.Printf("Creating %s from %s.\n", output, input)
	o, err := t.Render(data)
	if err != nil {
		return err
	}

	f.WriteString(o)
	f.Close()

	return nil
}
