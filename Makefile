NAME = animenamer
BINARY = animenamer

GO_FLAGS = #-v -race
GO_LDFLAGS = -ldflags "\
 -X github.com/tnextday/animenamer/cmd.AppVersion=`git describe --tags`\
 -X github.com/tnextday/animenamer/cmd.BuildTime=`date '+%Y-%m-%d_%H:%M:%S'`\
 -X github.com/tnextday/animenamer/cmd.DefaultTvDbApiKey=${TVDB_API_KEY}\
"
GO_VERSION = latest
GO_PROXY = https://goproxy.io


GOOS = `go env GOHOSTOS`
GOARCH = `go env GOHOSTARCH`

SOURCE_DIR = ./

all: local

.PHONY : local clean build linux-amd64 windows-amd64

clean:
	go clean -i $(GO_FLAGS) $(SOURCE_DIR)
	rm -f $(BINARY)
	rm -rf linux

fmt:
	goimports -w .

proxy:
	export GOPROXY=$(GO_PROXY)

mkdir:
	mkdir -p build/$(GOOS)-$(GOARCH)

build: mkdir
	go build $(GO_LDFLAGS) $(GO_FLAGS) -o build/$(GOOS)-$(GOARCH)/$(BINARY) $(SOURCE_DIR)

local: proxy build

darwin-amd64:
	mkdir -p build/$(GOOS)-$(GOARCH)
	GOOS=darwin GOARCH=amd64 go build $(GO_LDFLAGS) $(GO_FLAGS) -o build/darwin-amd64/$(BINARY) $(SOURCE_DIR)

linux-amd64:
	mkdir -p build/linux-amd64
	GOOS=linux GOARCH=amd64 go build $(GO_LDFLAGS) $(GO_FLAGS) -o build/linux-amd64/$(BINARY) $(SOURCE_DIR)

windows-amd64:
	mkdir -p build/windows-amd64
	GOOS=windows GOARCH=amd64 go build $(GO_LDFLAGS) $(GO_FLAGS) -o build/windows-amd64/$(BINARY).exe $(SOURCE_DIR)

release: proxy linux-amd64 windows-amd64
