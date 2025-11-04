# Build
FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY api/go.mod api/go.sum ./
RUN go mod download
COPY api/ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/quepc .

# Runtime
FROM alpine:3.20
RUN adduser -D -u 10001 appuser
WORKDIR /home/appuser
COPY --from=builder /bin/quepc /usr/local/bin/quepc
EXPOSE 8080
ENV APP_ENV=production APP_HOST=0.0.0.0 APP_PORT=8080
USER appuser
ENTRYPOINT ["/usr/local/bin/quepc"]