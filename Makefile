
generate:
	go generate ./...

test:
	go test -cover ./...

tidy:
	go mod tidy

server:
	go run main.go

dbinit:
	sqlite3 -echo -init sql/init.sql data/books.db  .quit

dbgen:
	sqlboiler  -o internal/models  sqlite3