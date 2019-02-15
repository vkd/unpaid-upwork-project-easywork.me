.PHONY: linux
linux:
	GOOS=linux GOARCH=amd64 go build -o server_linux
	upx -9 server_linux
