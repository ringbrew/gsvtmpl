PWD = $(shell pwd)

build-bin:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/{{.baseName}} cmd/{{.baseName}}.go
