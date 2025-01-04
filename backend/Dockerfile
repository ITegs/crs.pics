# +-----------------------------------------------------------------------------
# | Pull vendor dependencies
# +-----------------------------------------------------------------------------
FROM golang:1.23.4-alpine3.21 AS vendor

RUN mkdir /app && mkdir /build
WORKDIR /app

# Dependencies are installed in a separate layer, so they don't need to be
# reinstalled each time the source code changes.
COPY go.mod go.sum ./
RUN go mod download

# +-----------------------------------------------------------------------------
# | Base layer for the hot-reloading development images
# +-----------------------------------------------------------------------------
FROM vendor AS hot-reload

ENTRYPOINT [ "CompileDaemon" ]

# Install CompileDaemon, which will watch for changes in the source code and
# recompile the binaries on the fly.
RUN apk add --no-cache git
RUN go get github.com/githubnemo/CompileDaemon &&\
    go install github.com/githubnemo/CompileDaemon

# +-----------------------------------------------------------------------------
# | Build the binaries
# +-----------------------------------------------------------------------------
FROM vendor AS builder

# Build all binaries and make them available in the /build directory. The final
# image layer will copy them from there.
COPY . .
RUN go build -o /build main.go

# +-----------------------------------------------------------------------------
# | Assemble the final image, which doesn't need the source code or golang
# +-----------------------------------------------------------------------------
FROM alpine:3.18

# By setting /app as the WD, we will need to mount the config files under /app/config
RUN mkdir /app
WORKDIR /app

# Copy the compiled binaries from the builder image
COPY --from=builder /build/* /app/main

ENTRYPOINT ["/app/main"]
