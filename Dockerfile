FROM golang:1.22.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o shortener ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/shortener .

EXPOSE 8080

CMD ["./shortener", "-storage=postgres", "-pg_conn_str=postgres://postgres:postgres@host.docker.internal:5433/shortener?sslmode=disable"]
