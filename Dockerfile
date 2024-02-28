# Build Stage
FROM golang:1.20.14-alpine3.19 AS build

WORKDIR /app
COPY . .
RUN go build -o manu-auth .

# Run Stage
FROM alpine:3.19.1

WORKDIR /app
COPY --from=build /app/manu-auth .
EXPOSE 8080
CMD ["/app/manu-auth"]
