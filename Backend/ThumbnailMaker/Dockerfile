FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o thumbnail-maker ./cmd/main.go
# RUN go build -o video-service main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/thumbnail-maker .

COPY machinery.yaml .

RUN apk add --no-cache ffmpeg curl

EXPOSE 8081



CMD ["./thumbnail-maker"]