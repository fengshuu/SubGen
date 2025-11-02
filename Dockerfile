# syntax=docker/dockerfile:1

## Build stage
FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS build
WORKDIR /src
RUN apk add --no-cache git

# Pre-fetch dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/subgen ./main.go

## Runtime stage
FROM alpine:3.19
RUN apk add --no-cache ca-certificates && adduser -D -h /app appuser
WORKDIR /app
COPY --from=build /out/subgen ./subgen
COPY base_config.cache.yaml ./base_config.cache.yaml
EXPOSE 7081
RUN chown -R appuser:appuser /app
USER appuser
ENTRYPOINT ["./subgen"]
