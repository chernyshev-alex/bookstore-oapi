
generate:
	go generate ./...

test:
	go test -cover ./...

tidy:
	go mod tidy

server:
	go run cmd/srv/main.go