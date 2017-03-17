# Ricer

Generate your configuration files from templates.

The name comes from the popular linux past time, ricing. You can find more information [here](https://wiki.installgentoo.com/index.php/GNU/Linux_ricing).

The application thus aims to be only compatible with GNU/Linux, but can probably be used on macOS too.

## Installation

### Simple

  1. Download https://github.com/kvannotten/ricer/releases/download/v0.2/ricer_x64.tar.gz
  2. `tar xzvf ricer_x64.tar.gz`
  3. cp ricer/ricer /usr/bin/
  4. mkdir -p ~/.config/ricer/{plugins,templates}
  4. cp ricer/*.so ~/.config/ricer/plugins/

#### Using git

  1. `git clone https://github.com/kvannotten/ricer ~/go/src/github.com/kvannotten/ricer`
  2. `cd ~/go/src/github.com/kvannotten/ricer`
  3. `make # The makefile will compile and install the plugins in ~/.config/ricer/plugins`
  4. `make install`

## Basic usage

  1. Create a `XDG_HOME_CONFIG/ricer/templates` folder
  2. Add a configuration file `XDG_HOME_CONFIG/ricer/config.yml`
  3. Add templates to use
  4. Run ricer, it will automatically parse all config files and output them
  
Usually `XDG_HOME_CONFIG` is ~/.config
  
## Advanced

You can invoke ricer with an optional `env` parameter, that will inject the environment into your vars. For example, invoking ricer with: `ricer -env laptop` will add a `laptop: true` entry in yours vars section (at runtime).

This is handy if you want to consolidate all your templates and want to have certain parts only available depending on the enviroment available.

The default enviroment is "default".


## Config file

The config file has the following structure:

```
---
mytemplatename:
  input: /path/to/your/template # optional
  disabled: false # optional, so you can disable the rendering of this template
  output: /path/to/write/output/file
  engine: go_template # optional
  vars:
    Some: template
    Variables: you
    Want: to
    Use: !
```

Please note that input is optional, it will default to `XDG_HOME_CONFIG/ricer/templates/mytemplatename.tmpl`

Engine is also optional and will default to `go_template`. You can easily add your own plugins, but we provide support for `go_template` for the default golang templating language, and for `mustache`, the mustache handlebars templating engine.

## Templates

Using our configuration example from above, we should create the following file: `/path/to/your/template`, which could contain the following:


**Note:** These examples use the go [html/template](https://golang.org/pkg/html/template/) package for its templates. You can also choose to use the [mustache](https://mustache.github.io/) templating engine by specifying that engine in your template's configuration

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

## Plugins

Ricer exposes a plugin system that allows developers to define their own templating if they are unwilling to use the golang templating or the mustache templating. The plugin requires an `Execute` function to be defined:

```
func Execute(input, output string, data interface{})
```
The plugin has to be compiled with the `-buildmode=plugin` option, and has to have the `.so` extension. The name of the plugin is free for you to define, but you have to use the exact same name in the config file under the `engine` key, for it to work. Check the `plugins/` folder for examples.

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

### Tmuxinator configurations

I use tmuxinator to set up tmux sessions and do automatic layouting etc. A lot of those files look the same with a few changes. Ricer is excellent to automate this.

Consider this template:

```
# ~/.tmuxinator/{{name}}.yml

name: {{name}}
root: {{path}}

# Optional tmux socket
socket_name: {{name}}

windows:
    - editor:
      layout: main-vertical
      panes:
          - nvim
  - data:
      layout: tiled
      panes:
        - bundle exec rails db
        - bundle exec rails c
        - 
  - server:
      layout: main-vertical
      panes:
        - bundle exec rails s
```

We can then define several configuration entries that will output multiple files usable with tmuxinator.

```
projectone:
  input: /home/kristof/.config/ricer/templates/tmuxinator.tmpl
  output: /home/kristof/.tmuxinator/projectone.yml
  engine: mustache
  vars:
    name: projectone
    path: /path/to/projectone
projecttwo:
  input: /home/kristof/.config/ricer/templates/tmuxinator.tmpl
  output: /home/kristof/.tmuxinator/projecttwo.yml
  engine: mustache
  vars:
    name: projecttwo
    path: /path/to/projecttwo
```
