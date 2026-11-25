# Stage 1: Build static Go binary
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/api ./cmd/api

# Stage 2: Minimal runtime image
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata \
  && addgroup -S appgroup \
  && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/api ./api

RUN mkdir -p /app/data && chown -R appuser:appgroup /app/data

ENV PORT=8080
ENV ENV=production
ENV DB_PATH=/app/data/car_maintenance.db

USER appuser

EXPOSE 8080
VOLUME ["/app/data"]

CMD ["./api"]
