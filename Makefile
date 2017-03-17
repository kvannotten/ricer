VERSION="0.2"
BUILD_DATE="$(shell date -u +%Y-%m-%d_%H:%M)"
UPX := $(shell command -v upxfd 2> /dev/null)

all: ricer go_template.so mustache.so

go_template.so:
	go build -ldflags "-s" -buildmode=plugin -o go_template.so plugins/go_template.go

mustache.so:
	go build -ldflags "-s" -buildmode=plugin -o mustache.so plugins/mustache.go

ricer:
	go build -ldflags "-s -X main.version=$(VERSION) -X main.buildDate=$(BUILD_DATE)"
ifdef UPX
	upx ricer
endif

clean:
	rm -rf go_template.so mustache.so ricer

install:
	sudo cp ricer /usr/bin/ricer
	mkdir -p ~/.config/ricer/plugins
	cp *.so ~/.config/ricer/plugins/
