FROM golang:1.18-alpine AS golang-builder

ARG NameFile

RUN apk add bash && ls -a
COPY tmp/$NameFile.go /go/bin/$NameFile.go
RUN mv /go/bin/${NameFile}.go /go/bin/main.go
WORKDIR "/go/bin"
RUN go build main.go



CMD [$NameFile]
FROM scratch


ARG NameFile

COPY --from=golang-builder /go/bin/main /go/bin/main
WORKDIR "/go/bin"
CMD ["/go/bin/main"]
