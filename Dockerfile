# Stage 1: Modules caching
FROM golang:1.24 AS modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Stage 2: Build our application code
FROM golang:1.24 AS builder
ARG VERSION="v0.0.0-docker"
COPY --from=modules /go/pkg /go/pkg
COPY . /workdir
WORKDIR /workdir
RUN go build -ldflags "-X main.version=$VERSION" -o /bin/trustd ./cmd/trustd

# Stage 3: Final assembly
FROM gcr.io/distroless/base-debian12 AS trustd
COPY --from=builder /bin/trustd /
ENTRYPOINT ["/trustd"]