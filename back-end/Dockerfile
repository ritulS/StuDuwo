FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o api

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api .
COPY .env .
EXPOSE 5000
CMD ["./api"]
