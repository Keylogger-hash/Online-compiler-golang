FROM golang:1.18 as golang-builder
RUN mkdir -p "/go/src/sandbox/"
COPY sandbox-play/sandbox-play.go /go/src/sandbox/sandbox-play.go
 
WORKDIR "/go/src/sandbox"
RUN go mod init && go mod tidy && go build sandbox-play.go && chmod +x sandbox-play






FROM busybox:glibc

COPY --from=golang-builder /go/src/sandbox/sandbox-play .

CMD ["./sandbox-play"]
