# syntax=docker/dockerfile:1

## Build stage
FROM golang:1.21-alpine AS build
WORKDIR /src
RUN apk add --no-cache git

# Pre-fetch dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/subgen ./main.go

## Runtime stage
FROM alpine:3.19
RUN apk add --no-cache ca-certificates && adduser -D -h /app appuser
WORKDIR /app
COPY --from=build /out/subgen ./subgen
COPY base_config.cache.yaml ./base_config.cache.yaml
EXPOSE 8080
RUN chown -R appuser:appuser /app
USER appuser
ENTRYPOINT ["./subgen"]
