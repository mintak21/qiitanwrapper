FROM golang:1.14.7-buster as builder
COPY . /app
WORKDIR /app/cmd/api
RUN GO111MODULE=on go mod download && GOOS=linux GOARCH=amd64 go build -o /go/bin/qiitawrapper

# hadolint ignore=DL3007
FROM gcr.io/distroless/base-debian10:latest
COPY --from=builder /go/bin/qiitawrapper /bin/qiitawrapper
ENTRYPOINT [ "/bin/qiitawrapper" ]
