FROM golang:1.25.6-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 go build -o /out/api ./cmd/api

FROM golang:1.25.6-alpine AS dev
RUN apk add --no-cache make
RUN go install github.com/air-verse/air@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest
RUN adduser -D -u 1000 saleem
WORKDIR /app
USER saleem
EXPOSE 8080
CMD ["air", "-c", ".air.toml"]

FROM alpine:3.20 AS runtime
RUN adduser -D -u 10001 app
WORKDIR /app
COPY --from=builder /out/api /app/api
USER app
EXPOSE 8080
CMD ["/app/api"]