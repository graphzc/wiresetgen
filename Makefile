get-dev-alias:
	@echo 'alias wiresetgen="go run cmd/wiresetgen/main.go"'

build:
	@echo 'Building the project...'
	@go build -o wiresetgen cmd/wiresetgen/main.go
	@echo 'Build complete!'

build-and-install: build
	@echo 'Installing the project...'
	mv wiresetgen ~/go/bin
	@echo 'Installation complete!'

fmt:
	go fmt ./...