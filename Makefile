.PHONY: build_on_linux build_on_mac

build_on_linux:
	docker build -f ./deployment/dockerfile/qiitawrapper/Dockerfile -t linux_qiitawrapper .

build_on_mac:
	cd ./cmd/qiitawrapper; \
	GO111MODULE=on go mod download; \
	GOOS=darwin GOARCH=amd64 go build -o ../../qiitawrapper
