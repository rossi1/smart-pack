FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY app.env /app/app.env

RUN go build -o /bin/smart-pack

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /bin/smart-pack /app/smart-pack
COPY --from=builder /app/docker-entrypoint.sh /app/docker-entrypoint.sh
COPY --from=builder /app/app.env /app/app.env
COPY --from=builder /app/resources /app/resources
COPY --from=builder /app/api/ /app/api/

RUN chmod +x /app/docker-entrypoint.sh

ENTRYPOINT ["/app/docker-entrypoint.sh", "api"]
