FROM golang:1.18-alpine3.16 as base

WORKDIR /work
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY . .
RUN go mod download

FROM base as fish

RUN rm -rf pkg/proto/nrm
COPY protos/fish pkg/proto/nrm
RUN go build cmd/client/main.go
ENTRYPOINT ["go", "run", "cmd/client/main.go"]

# Default
FROM base

RUN go build cmd/client/main.go
ENTRYPOINT ["go", "run", "cmd/client/main.go"]
