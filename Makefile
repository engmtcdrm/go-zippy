.PHONY:	 test testv testcover testcoverall

test:
	@go test -timeout 30s ./...

testv:
	@go test -timeout 30s -v ./..

testcover:
	@go test -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html

testcoverall:
	@go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html
