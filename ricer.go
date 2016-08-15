/*
Ricer generates configuration files based on templates.
Copyright (C) 2016  Kristof Vannotten

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"os/user"
	"path"
	"sync"

	"github.com/spf13/viper"
)

var (
	version   string
	buildDate string

	config = flag.String("c", "", "(optional) the configuration file to use")
)

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

	var throttle = make(chan int, 4)
	var wg sync.WaitGroup

	for _, tmpl := range viper.AllKeys() {
		throttle <- 1
		wg.Add(1)

		go func(tmpl string, wg *sync.WaitGroup, throttle chan int) {
			defer wg.Done()
			if err := handleTemplate(tmpl); err != nil {
				fmt.Println(err)
			}
			<-throttle
		}(tmpl, &wg, throttle)
	}

	wg.Wait()
}

func parseConfiguration() error {
	if *config == "" {
		viper.SetConfigName("config")
		configHome, err := configHomeDirectory()
		if err != nil {
			return err
		}
		viper.AddConfigPath(configHome)
	} else {
		viper.SetConfigFile(*config)
	}

	return viper.ReadInConfig()
}

func configHomeDirectory() (string, error) {
	configHome := os.Getenv("XDG_CONFIG_HOME")

	if configHome == "" {
		homeDir := os.Getenv("HOME")
		if homeDir == "" {
			usr, err := user.Current()
			if err != nil {
				return "", err
			}
			homeDir = usr.HomeDir
		}
		configHome = path.Join(homeDir, ".config")
	}

	return path.Join(configHome, "ricer"), nil
}

func templatePath(tmpl string) (string, error) {
	config, err := configHomeDirectory()
	if err != nil {
		return "", err
	}

	return path.Join(config, "templates", fmt.Sprintf("%s.tmpl", tmpl)), nil
}

func handleTemplate(tmpl string) error {
	var t *template.Template
	var err error

	// parse the template
	file := viper.GetString(fmt.Sprintf("%s.input", tmpl))
	if file == "" {
		if file, err = templatePath(tmpl); err != nil {
			return err
		}
	}

	if t, err = template.ParseFiles(file); err != nil {
		return fmt.Errorf("Could not parse template %s", file)
	}

	// get configuration details
	m := viper.GetStringMap(fmt.Sprintf("%s.vars", tmpl))
	outputFile := viper.GetString(fmt.Sprintf("%s.output", tmpl))

	// check if an output file is given
	if outputFile == "" {
		return fmt.Errorf("You have to define an output for template %s", tmpl)
	}

	// create the path of the output file
	path := path.Dir(outputFile)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fmt.Errorf("[1] Could not create %s for template %s", outputFile, tmpl)
	}

	// write the output file
	var f *os.File
	if f, err = os.Create(outputFile); err != nil {
		return fmt.Errorf("[2] Could not create %s for template %s", outputFile, tmpl)
	}

	fmt.Printf("Creating %s from template %s.\n", outputFile, file)
	t.Execute(f, m)
	f.Close()

	return nil
}
