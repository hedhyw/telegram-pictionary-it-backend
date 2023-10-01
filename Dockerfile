ARG GOLANG_DOCKER_TAG=1.21.1-alpine3.18
ARG ALPINE_DOCKER_TAG=3.18

FROM golang:$GOLANG_DOCKER_TAG as builder

RUN apk update && apk upgrade && apk add --no-cache make curl

WORKDIR /build
COPY . .

RUN make build OUT_BUILD=/build/bin/server

FROM alpine:$ALPINE_DOCKER_TAG

WORKDIR /app

COPY --from=builder /build/bin/server /app/server

ENTRYPOINT [ "/app/server" ]
