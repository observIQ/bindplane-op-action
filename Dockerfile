FROM golang:1.26.5-alpine AS builder
WORKDIR /app
COPY . .
WORKDIR /app/cmd/action
RUN CGO_ENABLED=0 go build -o /entrypoint

FROM alpine:3.24.1
RUN apk add --no-cache ca-certificates
COPY --from=builder /entrypoint /entrypoint
ENTRYPOINT ["/entrypoint"]
