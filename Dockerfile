FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o avito-shop ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/avito-shop .

EXPOSE 8080

CMD ["./avito-shop"]