FROM golang:1.22 as build
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o service ./cmd/auction

# FROM scratch
FROM alpine:latest
RUN apk add --no-cache bash
WORKDIR /app
COPY ./cmd/auction/.env .
COPY --from=build /app/service .
ENTRYPOINT ["./service"]