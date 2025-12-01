FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o video-service ./cmd/main.go
# RUN go build -o video-service main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/video-service .
COPY machinery.yaml .
EXPOSE 8081

CMD ["./video-service"]