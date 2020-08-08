.PHONY: build_on_linux build_on_mac

build_on_linux:
	docker build -f ./deployment/dockerfile/api/Dockerfile -t linux_qiitawrapper .

build_on_mac:
	cd ./cmd/api; \
	GO111MODULE=on go mod download; \
	GOOS=darwin GOARCH=amd64 go build -o ../../qiitawrapper
