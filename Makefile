build:
	go build -o bin/monkey-language-server cmd/monkey-lsp/main.go

run:
	@go run cmd/monkey-lsp/main.go
