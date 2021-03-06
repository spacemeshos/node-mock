# Inspired by https://container-solutions.com/faster-builds-in-docker-with-go-1-11/
# Base build image
FROM golang:1.13-alpine AS build_base
RUN apk add bash make git curl unzip rsync libc6-compat gcc musl-dev
WORKDIR /go/src/github.com/spacemeshos/node-mock

# Force the go compiler to use modules
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

RUN go get github.com/golang/snappy@v0.0.1

# This image builds the node-mock
FROM build_base AS server_builder
# Here we copy the rest of the source code
COPY . .

# And compile the project
RUN go build

#In this last stage, we start from a fresh Alpine image, to reduce the image size and not ship the Go compiler in our production artifacts.
FROM alpine AS node-mock

# Finally we copy the statically compiled Go binary.
COPY --from=server_builder /go/src/github.com/spacemeshos/node-mock/node-mock /bin/node-mock

ENTRYPOINT ["/bin/node-mock"]
# profiling port
EXPOSE 6060

# gRPC port
EXPOSE 9990
