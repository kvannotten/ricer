package main

import (
	"fmt"
	"html/template"
	"os"
)

// Execute the template
func Execute(input, output string, data interface{}) error {
	var t *template.Template
	var err error

	if t, err = template.ParseFiles(input); err != nil {
		return fmt.Errorf("Could not parse template %s", input)
	}

	// write the output file
	var f *os.File
	if f, err = os.Create(output); err != nil {
		return fmt.Errorf("[2] Could not create %s from %s", output, input)
	}

	fmt.Printf("Creating %s from %s.\n", output, input)
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	f.Close()

	return nil
}
