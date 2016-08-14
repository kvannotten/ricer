package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var version string
var buildDate string

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ricer version %s\n", version)
		fmt.Fprintf(os.Stderr, "Built %s\n", buildDate)
		fmt.Fprintln(os.Stderr, "usage:")
		flag.PrintDefaults()
	}
	flag.Parse()

	if err := parseConfiguration(); err != nil {
		panic(err)
	}

	files, _ := filepath.Glob("*.tmpl")
	for _, file := range files {
		var t *template.Template
		var err error

		// parse the template
		if t, err = template.ParseFiles(file); err != nil {
			fmt.Printf("Could not parse template %s\n", file)
			continue
		}

		// get the template name based on the filename
		tmplName := strings.TrimSuffix(file, filepath.Ext(file))

		// get configuration details
		m := viper.GetStringMapString(fmt.Sprintf("%s.vars", tmplName))
		outputFile := viper.GetString(fmt.Sprintf("%s.output", tmplName))

		// check if an output file is given
		if outputFile == "" {
			fmt.Printf("You have to define an output for template %s\n", tmplName)
			continue
		}

		// create the path of the output file
		path := path.Dir(outputFile)
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			fmt.Printf("[1] Could not create %s for template %s\n", outputFile, tmplName)
			continue
		}

		// write the output file
		if f, err := os.Create(outputFile); err != nil {
			fmt.Printf("[2] Could not create %s for template %s\n", outputFile, tmplName)
			continue
		} else {
			t.Execute(f, m)
			f.Close()
		}
	}
}

func parseConfiguration() error {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	return viper.ReadInConfig()
}
