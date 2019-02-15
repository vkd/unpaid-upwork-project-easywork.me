run:
	go run cmd/tracker/main.go

vendor:
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor
