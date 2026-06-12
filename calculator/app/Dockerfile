
FROM golang:1.26 AS builder

WORKDIR /app

COPY . .

RUN go build -o calculator .

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/calculator /app/calculator

EXPOSE 8080

CMD ["/app/calculator"]

