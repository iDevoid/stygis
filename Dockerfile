FROM golang:1.14-alpine AS build_base

RUN apk --no-cache add tzdata
RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /tmp/go-app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Unit tests
RUN CGO_ENABLED=0 go test ./... -v

# Build the Go app
RUN go build -o ./bin/stygis cmd/main.go

# Start fresh from a smaller image
FROM alpine:3.9 
RUN apk add ca-certificates

COPY --from=build_base /tmp/go-app/bin/stygis /app/stygis
COPY --from=build_base /tmp/go-app/static /static
COPY --from=build_base /usr/share/zoneinfo /usr/share/zoneinfo

ENV TZ=Asia/Jakarta
# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
CMD ["/app/stygis"]