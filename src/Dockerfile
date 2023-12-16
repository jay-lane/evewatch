FROM golang:1.14 as builder
WORKDIR /go/src/github.com/jay-lane/evewatch
COPY . .
RUN make build

FROM quay.io/prometheus/busybox:latest AS app

COPY --from=builder /go/src/github.com/jay-lane/evewatch /bin/evewatch

EXPOSE 9999
ENTRYPOINT ["/bin/evewatch"]