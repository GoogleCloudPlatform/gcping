release:
	GOOS=windows GOARCH=amd64 go build -o ./bin/gcping_windows_amd64
	GOOS=linux GOARCH=amd64 go build -o ./bin/gcping_linux_amd64
	GOOS=darwin GOARCH=amd64 go build -o ./bin/gcping_darwin_amd64

push:
	gsutil cp bin/* gs://gcping-release
