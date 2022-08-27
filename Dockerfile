FROM golang:1.18-alpine3.16

WORKDIR /work
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY . .
RUN go mod download
RUN go build cmd/client/main.go

ENTRYPOINT ["go", "run", "cmd/client/main.go"]
