.PHONY: build_on_linux build_on_mac

build_on_linux:
	cd ./cmd/qiitawrapper; \
	GO111MODULE=on go mod download; \
	GOOS=linux GOARCH=amd64 go build -o ../../qiitawrapper

build_on_mac:
	cd ./cmd/qiitawrapper; \
	GO111MODULE=on go mod download; \
	GOOS=darwin GOARCH=amd64 go build -o ../../qiitawrapper

