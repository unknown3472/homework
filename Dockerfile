# Stage 1: Build the Go application
FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /taskapp cmd/main.go

# Stage 2: Create a smaller image to run the application
FROM alpine:latest

WORKDIR /

COPY --from=builder /taskapp /taskapp

EXPOSE 8080

CMD ["/taskapp"]
