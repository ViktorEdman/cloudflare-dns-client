build: main.go
	GOOS=windows GOARCH=amd64 go build -o bin/cfdns-amd64-windows.exe .
	GOOS=linux GOARCH=amd64 go build -o bin/cfdns-amd64-linux .
	GOOS=linux GOARCH=arm go build -o bin/cfdns-arm-linux .
	GOOS=linux GOARCH=arm64 go build -o bin/cfdns-arm64-linux .
