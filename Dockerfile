FROM golang:1.18-alpine as builder
ENV CGO_ENABLED=0
WORKDIR /go/src/
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags '-w -s' -v -o /usr/local/bin/ksops ./

FROM alpine:latest
COPY --from=builder /usr/local/bin/ksops /usr/local/bin/ksops
ENTRYPOINT ["ksops"]
