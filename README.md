# Ricer

Generate your configuration files from templates.

## Installation

#### Using go get

`go get github.com/kvannotten/ricer`

#### Using git

  1. `git clone https://github.com/kvannotten/ricer ~/go/src/github.com/kvannotten/ricer`
  2. `cd ~/go/src/github.com/kvannotten/ricer`
  3. `make`
  4. `cp ricer SOMEWHEREINYOURPATH`

## Usage

  1. Create a `XDG_HOME_CONFIG/ricer/templates` folder
  2. Add a configuration file `XDG_HOME_CONFIG/ricer/config.yml`
  3. Add templates to use
  4. Run ricer, it will automatically parse all config files and output them

Usually `XDG_HOME_CONFIG` is ~/.config

## Config file

The config file has the following structure:

```
---
mytemplatename:
  output: /path/to/write/output/file
  vars:
    Some: template
    Variables: you
    Want: to
    Use: !
```

## Templates

You have to have a configuration entry for every template you put in the templates directory.

Templates' extensions should be `tmpl`, thus you can disable templates by using a different extension.

Using our configuration example from above, we should create the following file: `XDG_HOME_CONFIG/ricer/templates/mytemplatename.tmpl`, which could contain the following:

```
let my_option = {{.Some}}
let option2 = {{.Variables}}
{{.Want}} = {{.Use}}
```

Using the configuration file from above, this would result in the following file:

```
let my_option = template
let option2 = you
to = !
```
