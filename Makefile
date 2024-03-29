NAME = animenamer
BINARY = animenamer

GO_VERSION = latest
GO_PROXY = https://goproxy.io
GIT_TAG = `git describe --tags`
GOOS = `go env GOHOSTOS`
GOARCH = `go env GOHOSTARCH`
TVDB_APIKEY = Z5SC1ZD07NNS8TDC
TMDB_APIKEY = 4219e299c89411838049ab0dab19ebd5

GO_FLAGS = #-v -race
GO_LDFLAGS = -ldflags "\
 -X github.com/tnextday/animenamer/cmd.AppVersion=$(GIT_TAG)\
 -X github.com/tnextday/animenamer/cmd.BuildTime=`date '+%Y-%m-%d_%H:%M:%S'`\
 -X github.com/tnextday/animenamer/cmd.DefaultTVDBApiKey=${TVDB_APIKEY}\
 -X github.com/tnextday/animenamer/cmd.DefaultTMDBApiKey=${TMDB_APIKEY}\
"

SOURCE_DIR = ./

all: local

.PHONY : local clean build release

clean:
	go clean -i $(GO_FLAGS) $(SOURCE_DIR)
	rm -f $(BINARY)
	rm -rf build
	rm -rf _tests

fmt:
	goimports -w .

proxy:
	export GOPROXY=$(GO_PROXY)

build:
	mkdir -p build/$(GOOS)-$(GOARCH)
	go build $(GO_LDFLAGS) $(GO_FLAGS) -o build/$(GOOS)-$(GOARCH)/$(BINARY) $(SOURCE_DIR)
	ln -sf $(GOOS)-$(GOARCH)/$(BINARY) build/$(BINARY)

local: proxy build

darwin-amd64:
	mkdir -p build/darwin-amd64
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(GO_LDFLAGS) $(GO_FLAGS) -o build/darwin-amd64/$(BINARY) $(SOURCE_DIR)
	cd build/darwin-amd64 && zip ../releases/animenamer_$(GIT_TAG)_darwin_amd64.zip animenamer

linux-amd64:
	mkdir -p build/linux-amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(GO_LDFLAGS) $(GO_FLAGS) -o build/linux-amd64/$(BINARY) $(SOURCE_DIR)
	cd build/linux-amd64 && zip ../releases/animenamer_$(GIT_TAG)_linux_amd64.zip animenamer

linux-386:
	mkdir -p build/linux-386
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build $(GO_LDFLAGS) $(GO_FLAGS) -o build/linux-386/$(BINARY) $(SOURCE_DIR)
	cd build/linux-386 && zip ../releases/animenamer_$(GIT_TAG)_linux_386.zip animenamer

linux-arm:
	mkdir -p build/linux-arm
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build $(GO_LDFLAGS) $(GO_FLAGS) -o build/linux-arm/$(BINARY) $(SOURCE_DIR)
	cd build/linux-arm && zip ../releases/animenamer_$(GIT_TAG)_linux_arm.zip animenamer

linux-arm64:
	mkdir -p build/linux-arm64
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(GO_LDFLAGS) $(GO_FLAGS) -o build/linux-arm64/$(BINARY) $(SOURCE_DIR)
	cd build/linux-arm64 && zip ../releases/animenamer_$(GIT_TAG)_linux_arm64.zip animenamer


windows-amd64:
	mkdir -p build/windows-amd64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(GO_LDFLAGS) $(GO_FLAGS) -o build/windows-amd64/$(BINARY).exe $(SOURCE_DIR)
	cd build/windows-amd64 && zip ../releases/animenamer_$(GIT_TAG)_windows_amd64.zip animenamer.exe

windows-386:
	mkdir -p build/windows-386
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build $(GO_LDFLAGS) $(GO_FLAGS) -o build/windows-386/$(BINARY).exe $(SOURCE_DIR)
	cd build/windows-386 && zip ../releases/animenamer_$(GIT_TAG)_windows_386.zip animenamer.exe

release-dir:
	mkdir -p build/releases

release: proxy release-dir darwin-amd64 linux-amd64 linux-386 linux-arm linux-arm64 windows-amd64 windows-386
