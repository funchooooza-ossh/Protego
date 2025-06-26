FROM golang:1.24 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o protego ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o cli-tool ./cmd/cli

FROM debian:bullseye-slim

RUN apt-get update \
 && apt-get install -y --fix-missing ca-certificates \
 && rm -rf /var/lib/apt/lists/*


WORKDIR /app
COPY --from=builder /app/protego .
COPY --from=builder /app/cli-tool .
COPY .env .
COPY entrypoint.sh .
COPY migrations ./migrations
RUN chmod +x entrypoint.sh

EXPOSE 8080
ENTRYPOINT ["/bin/sh", "./entrypoint.sh"]
