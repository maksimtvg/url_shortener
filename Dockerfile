FROM golang:1.18 AS builder
WORKDIR /app
COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags='-w -s' -o /go/bin/app ./cmd/main.go

FROM alpine:latest
COPY --from=builder /go/bin/app /go/bin/app
ENTRYPOINT ["/go/bin/app"]