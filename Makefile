build:
	@echo 'Building the project...'
	@go build -o wiresetgen cmd/wiresetgen/main.go
	@echo 'Build complete!'

install:
	@echo 'Installing the project...'
	go install ./cmd/wiresetgen
	@echo 'Install complete!'
	@echo 'You can now run wiresetgen from anywhere!'

fmt:
	@echo 'Formatting the code...'
	go fmt ./...
	@echo 'Formatting complete!'

test:
	@echo 'Running tests...'
	go test ./...
	@echo 'Tests complete!'
