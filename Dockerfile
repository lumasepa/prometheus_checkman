FROM golang:1.10.3-alpine3.8 AS builder

WORKDIR /go/src/github.com/lumasepa/prometheus_checkman
ADD . .

RUN set -e; \
    apk --no-cache add ca-certificates; \
    apk --no-cache add --virtual build-deps git

RUN go get
RUN go build -o /bin/prometheus_checkman

FROM alpine:3.8
COPY --from=builder /bin/prometheus_checkman /bin/prometheus_checkman
ENTRYPOINT [ "/bin/prometheus_checkman" ]
