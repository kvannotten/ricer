VERSION="0.2"
BUILD_DATE="$(shell date -u +%Y-%m-%d_%H:%M)"

all:
	go build -ldflags "-s -X main.version=$(VERSION) -X main.buildDate=$(BUILD_DATE)"
	mkdir -p ~/.config/ricer/plugins
	go build -ldflags "-s" -buildmode=plugin -o ~/.config/ricer/plugins/go_template.so plugins/go_template.go
	go build -ldflags "-s" -buildmode=plugin -o ~/.config/ricer/plugins/mustache.so plugins/mustache.go
	upx ricer
