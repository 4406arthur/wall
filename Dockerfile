# build stage
FROM golang:1.17.2-alpine AS build-env
ADD . /src
RUN \
    cd /src && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -o /wall

# final stage
FROM alpine:3.15
COPY --from=build-env /wall /usr/local/bin

RUN \
    apk update \
    && apk add ca-certificates \
    && rm -rf /var/cache/apk/*

EXPOSE 80 443

ENTRYPOINT ["wall"]
