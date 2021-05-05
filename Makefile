.PHONY: build

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./Handler ./api/main.go
	zip ./Handler.zip ./Handler

build-local:
	go build -o ./Handler ./api/main.go