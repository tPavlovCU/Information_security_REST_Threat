# Этап сборки
FROM golang:1.27.3-alpine AS builder

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build -o /threat-api ./cmd/api


FROM alpine:3.22

WORKDIR /app


RUN apk --no-cache add ca-certificates

COPY --from=builder /threat-api ./threat-api

EXPOSE 8080

CMD ["./threat-api"]