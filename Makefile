VERSION="0.2"
BUILD_DATE := "$(shell date -u +%Y-%m-%d_%H:%M)"
UPX := $(shell command -v upx 2> /dev/null)
GOCMD   := go
GOBUILD := $(GOCMD) build

all: ricer go_template.so mustache.so

%.so: plugins/%.go
	echo "FOO"
	$(GOBUILD) -ldflags "-s" -buildmode=plugin -o $@ $<

ricer:
	$(GOBUILD) -ldflags "-s -X main.version=$(VERSION) -X main.buildDate=$(BUILD_DATE)"
ifdef UPX
	$(UPX) ricer
endif

clean:
	rm -rf go_template.so mustache.so ricer

install:
	sudo cp ricer /usr/bin/ricer
	mkdir -p ~/.config/ricer/plugins
	cp *.so ~/.config/ricer/plugins/
