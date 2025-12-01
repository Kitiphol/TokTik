FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o user-service ./cmd/main.go
# RUN go build -o user-service main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/user-service .
EXPOSE 8081

CMD ["./user-service"]