FROM golang:1.22-alpine as builder
WORKDIR /app
COPY . .
WORKDIR /app/cmd/action
RUN CGO_ENABLED=0 go build -o /entrypoint

FROM alpine:3.10
RUN apk add --no-cache git
COPY --from=builder /entrypoint /entrypoint
ENTRYPOINT ["/entrypoint"]
