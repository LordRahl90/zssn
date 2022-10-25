start:
	@go run ./cmd/

test:
	@go test ./... --cover

build-image:
	@docker build -t gcr.io/neurons-be-test/zssn:latest .

push-image:
	@docker push gcr.io/neurons-be-test/zssn:latest

bi: build-image
pi: push-image