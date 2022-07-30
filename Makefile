
generate:
	go generate ./...

test:
	go test -cover ./...

tidy:
	go mod tidy

server:
	go run cmd/srv/main.go

dbcreate:
	sqlite3 -echo -init sql/init.sql data/books.db

dbboil:
	sqlboiler  -o internal/models  sqlite3