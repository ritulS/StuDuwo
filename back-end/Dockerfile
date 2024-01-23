FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o api



FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api .

EXPOSE 3000

CMD ["./api"]
