test: 
	go test ./...

build: 
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/redis-lite cmd/redis-lite/main.go

run:
	go run cmd/redis-lite/main.go