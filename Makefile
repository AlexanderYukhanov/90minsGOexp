.PHONY: swagger
swagger:
	go generate swagger/gen.go

build: bin/server
	go build -o bin/server server/cmd/experimental-server/main.go

.PHONY: start
start:
	bin/server --port 8080