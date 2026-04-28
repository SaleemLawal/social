FROM golang:1.25.6-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 go build -o /out/api ./cmd/api

FROM alpine:3.20 AS runtime
RUN adduser -D -u 10001 app
WORKDIR /app
COPY --from=builder /out/api /app/api
USER app
EXPOSE 8080
CMD ["/app/api"]