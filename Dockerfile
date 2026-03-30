FROM golang:1.25-alpine AS builder

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=off

COPY . .

RUN go build -o /app/server ./cmd

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/server /app/server
COPY --from=builder /app/internal/config /app/internal/config
COPY --from=builder /app/db /app/db

EXPOSE 8080

ENV CONFIG_PATH=/app/internal/config/config.docker.yaml

CMD ["/app/server"]
