FROM golang:1.26.10-alpine AS builder
RUN apk upgrade --no-cache ca-certificates
WORKDIR /app
COPY . .
WORKDIR /app/cmd/action
RUN CGO_ENABLED=0 go build -o /entrypoint

FROM alpine:3.24.3
RUN apk add --no-cache --upgrade ca-certificates
COPY --from=builder /entrypoint /entrypoint
ENTRYPOINT ["/entrypoint"]
