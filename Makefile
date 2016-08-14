VERSION="0.1"
BUILD_DATE="$(shell date -u +%Y-%m-%d_%H:%M)"

all:
	go build -ldflags "-s -X main.version=$(VERSION) -X main.buildDate=$(BUILD_DATE)"
	upx ricer
