ifeq ($(origin version), undefined)
	version := dev
endif

release:
	GOOS=windows GOARCH=amd64 go build -o ./out/gcping_windows_amd64_$(version)
	cp ./out/gcping_windows_amd64_$(version) ./out/gcping_windows_arm64_latest
	GOOS=linux GOARCH=amd64 go build -o ./out/gcping_linux_amd64_$(version)
	cp ./out/gcping_linux_amd64_$(version) ./out/gcping_linux_amd64_latest
	GOOS=darwin GOARCH=amd64 go build -o ./out/gcping_darwin_amd64_$(version)
	cp ./out/gcping_darwin_amd64_$(version) ./out/gcping_darwin_amd64_latest
	GOOS=darwin GOARCH=arm64 go build -o ./out/gcping_darwin_arm64_$(version)
	cp ./out/gcping_darwin_arm64_$(version) ./out/gcping_darwin_arm64_latest

push:
	gsutil cp out/* gs://gcping-release
