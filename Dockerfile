FROM golang:1.25.7-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /arena-ban-bin ./cmd/api/main.go

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /root/

COPY --from=builder /arena-ban-bin .

EXPOSE 3000

CMD ["./arena-ban-bin"]