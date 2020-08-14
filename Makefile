.PHONY: build_on_linux build_on_mac generate

build_on_linux:
	docker build -f ./deployment/dockerfile/api/Dockerfile -t linux_qiitawrapper .

build_on_mac:
	cd ./cmd/api; \
	GO111MODULE=on go mod download; \
	GOOS=darwin GOARCH=amd64 go build -o ../../qiitawrapper

generate:
	swagger generate server -a qiitawrapper -A qiitawrapper --exclude-main --strict-additional-properties -t gen -f ./swagger/swagger.yaml
