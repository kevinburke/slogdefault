test:
	go vet ./...
	go test -race ./...
