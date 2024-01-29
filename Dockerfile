# Container image that runs your code
FROM alpine:3.10

RUN apk add --no-cache bash curl git jq

COPY --chmod=0755 entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
