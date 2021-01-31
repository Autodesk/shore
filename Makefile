VERSION=devel
ARTIFACT_NAME=shore
BUILD_CMD=go build -o $(ARTIFACT_NAME) cmd/shore/shore.go
LD_FLAGS=-ldflags="-s -w -X 'main.Version=$(VERSION)'"

setup:
	$(shell export GOPRIVATE="github.com")
	$(shell export GOPROXY="https://:@***REMOVED***/***REMOVED***/gocenter/")
	go mod download
	go mod vendor

test:
	go test ./... -race -cover -v

build: build-osx build-linux build-win

build-osx:
	GOOS=darwin GOARCH=amd64 $(BUILD_CMD)

build-linux:
	GOOS=linux GOARCH=amd64 $(BUILD_CMD)

build-win:
	GOOS=windows GOARCH=amd64 $(BUILD_CMD)

build-release: build-release-osx build-release-linux build-release-win

build-release-osx:
	GOOS=darwin GOARCH=amd64 go build $(LD_FLAGS) -o $(ARTIFACT_NAME)-osx cmd/shore/shore.go

build-release-linux:
	GOOS=linux GOARCH=amd64 go build $(LD_FLAGS) -o $(ARTIFACT_NAME)-linux cmd/shore/shore.go

build-release-win:
	GOOS=windows GOARCH=amd64 go build $(LD_FLAGS) -o $(ARTIFACT_NAME)-win cmd/shore/shore.go

clean:
	rm -rf shore-*