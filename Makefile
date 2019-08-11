ifeq ($(origin version), undefined)
	version := latest
endif

release:
	GOOS=windows GOARCH=amd64 go build -o ./bin/gcping_windows_amd64_$(version)
	GOOS=linux GOARCH=amd64 go build -o ./bin/gcping_linux_amd64_$(version)
	GOOS=darwin GOARCH=amd64 go build -o ./bin/gcping_darwin_amd64_$(version)

push:
	gsutil cp bin/* gs://gcping-release
