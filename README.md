# Ricer

Generate your configuration files from templates.

The name comes from the popular linux past time, ricing. You can find more information [here](https://wiki.installgentoo.com/index.php/GNU/Linux_ricing).

The application thus aims to be only compatible with GNU/Linux, but can probably be used on macOS too.

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
  input: /path/to/your/template
  output: /path/to/write/output/file
  vars:
    Some: template
    Variables: you
    Want: to
    Use: !
```

Please note that input is optional, it will default to `XDG_HOME_CONFIG/ricer/templates/mytemplatename.tmpl`

## Templates

Using our configuration example from above, we should create the following file: `/path/to/your/template`, which could contain the following:


**Note:** Ricer uses the go [html/template](https://golang.org/pkg/html/template/) package for its templates

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

## Real life examples

#### i3 screens

I use i3 as a window manager on both my desktop and my laptop, on my laptop I use 3 screens and on my desktop 2. So I need a way to define these screens in i3. Ricer can help with that.

Conider this template:

```
# add screens
{{ range $screen, $value := .Screens }}set ${{ $screen }} {{ $value }}
{{ end }}
```

And this configuration snippet:

```
i3:
  input: /home/kristof/.config/ricer/templates/i3_desktop.tmpl
  output: /home/kristof/.i3/config
  vars:
    Term: /usr/local/bin/st
    Screens:
      screen1: DVI-I-3
      screen2: DVI-I-2
```

Runing `ricer` would result in the following i3 configuration output:

```
# add screens
set $screen1 DVI-I-3
set $screen2 DVI-I-2
```
