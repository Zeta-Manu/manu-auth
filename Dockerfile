# Build Stage
FROM golang:1.20.14-alpine3.19 AS build

WORKDIR /app
COPY . .
RUN go build -o manu-auth ./cmd/main.go

# Run Stage
FROM alpine:3.19.1

WORKDIR /app
COPY --from=build /app/manu-auth .
EXPOSE 8080

# Set GIN_MODE to release
ENV GIN_MODE=release

# Set APP env
ENV APP_HTTP_PORT=8080

CMD ["/app/manu-auth"]
