FROM golang:alpine AS server

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy the code into the container
COPY . .
RUN go mod download

# Build the application
RUN go build -o run -ldflags '-X main.buildTime=$(DATE)' cmd/gojastin/main.go

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/run .


# Build a small image
FROM scratch

COPY --from=server /dist/run /
# Command to run
ENTRYPOINT ["/run"]
