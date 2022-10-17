setup:
	go mod download
	go mod vendor
	go mod tidy

build:
	go build -o jt cmd/jt/main.go

release:
	$(shell curl -sL https://git.io/goreleaser | bash)

clean:
	rm -rf dist/
	rm -rf jt
