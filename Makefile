start:
	@go run ./cmd/

test:
	@go test ./... -v --cover