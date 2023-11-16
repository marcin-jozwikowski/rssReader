build:
	go build -ldflags "-s -w" reader.go

run:
	go run reader.go