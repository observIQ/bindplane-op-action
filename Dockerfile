FROM golang:1.26-alpine as builder
WORKDIR /app
COPY . .
WORKDIR /app/cmd/action
RUN CGO_ENABLED=0 go build -o /entrypoint

FROM alpine:3.10
RUN apk add --no-cache ca-certificates
COPY --from=builder /entrypoint /entrypoint
ENTRYPOINT ["/entrypoint"]
