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
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"plugin"
	"sync"

	"github.com/spf13/viper"
)

var (
	version   string
	buildDate string

	config      = flag.String("c", "", "(optional) the configuration file to use")
	environment = flag.String("env", "default", "(optional) the environment that will be injected in every vars entry")
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
		if viper.GetBool(fmt.Sprintf("%s.disabled", tmpl)) {
			continue
		}

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
	file, err := inputFile(tmpl)
	if err != nil {
		return err
	}

	// get configuration details
	m, err := templateVars(tmpl)
	if err != nil {
		return err
	}

	outputFile, err := outputFile(tmpl)
	if err != nil {
		return err
	}

	if err = createDirectoryForPath(outputFile); err != nil {
		return err
	}

	content, err := contentForTemplate(tmpl, file, m)
	if err != nil {
		return err
	}

	if !templateContentChanged(content, outputFile) {
		fmt.Printf("%s has not changed. Skipping...\n", outputFile)
		return nil
	}

	if err = writeTemplateContentToFile(content, outputFile); err != nil {
		return err
	}

	fmt.Printf("Created %s from %s\n", outputFile, file)

	return nil
}

func pluginPath(pluginName string) (string, error) {
	config, err := configHomeDirectory()
	if err != nil {
		return "", err
	}

	return path.Join(config, "plugins", fmt.Sprintf("%s.so", pluginName)), nil
}

func getTemplatingMethod(pluginName string) func(string, interface{}) ([]byte, error) {
	path, _ := pluginPath(pluginName)
	p, err := plugin.Open(path)
	if err != nil {
		panic(err)
	}

	execute, err := p.Lookup("Execute")
	if err != nil {
		panic(err)
	}

	return execute.(func(string, interface{}) ([]byte, error))
}

func inputFile(tmpl string) (string, error) {
	var err error

	file := viper.GetString(fmt.Sprintf("%s.input", tmpl))
	if file == "" {
		if file, err = templatePath(tmpl); err != nil {
			return "", err
		}
	}

	return file, nil
}

func templateVars(tmpl string) (map[string]interface{}, error) {
	m := viper.GetStringMap(fmt.Sprintf("%s.vars", tmpl))
	m[*environment] = true

	return m, nil
}

func outputFile(tmpl string) (string, error) {
	outputFile := viper.GetString(fmt.Sprintf("%s.output", tmpl))
	if outputFile == "" {
		return "", fmt.Errorf("No output file defined for template %s", tmpl)
	}

	return outputFile, nil
}

func createDirectoryForPath(file string) error {
	// create the path of the output file
	path := path.Dir(file)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fmt.Errorf("Failed to create path %s for %s.", path, file)
	}

	return nil
}

func contentForTemplate(tmpl, file string, m map[string]interface{}) ([]byte, error) {
	templateEngine := viper.GetString(fmt.Sprintf("%s.engine", tmpl))
	if templateEngine == "" {
		templateEngine = "go_template"
	}

	execute := getTemplatingMethod(templateEngine)
	return execute(file, m)
}

func writeTemplateContentToFile(content []byte, output string) error {
	var err error

	var f *os.File
	if f, err = os.Create(output); err != nil {
		return fmt.Errorf("[2] Could not create %s", output)
	}
	f.Write(content)
	f.Close()

	return nil
}

func templateContentChanged(content []byte, output string) bool {
	c, err := ioutil.ReadFile(output)
	if err != nil {
		return true
	}

	return !bytes.Equal(content, c)
}
