FROM golang:latest

EXPOSE 8080

ARG TARGETPLATFORM

LABEL maintainer="Dmitry Ponomarev <demdxx@gmail.com>"
LABEL service.name=rtb.sspserver
LABEL service.weight=1
LABEL service.port=8080
LABEL service.check.http=http://{{address}}/healthcheck
LABEL service.check.interval=5s
LABEL service.check.timeout=2s
LABEL service.public=false

ENV SERVER_HTTP_LISTEN=:8080
ENV SERVER_GRPC_LISTEN=tcp://:8081
ENV SERVER_PROFILE_MODE=net
ENV SERVER_PROFILE_LISTEN=8082
ENV ADSTORAGE_CONNECTION=fs:///data/

COPY .build/${TARGETPLATFORM}/sspserver /sspserver
COPY deploy/data /data
COPY public /public

ENTRYPOINT [ "/sspserver"]
