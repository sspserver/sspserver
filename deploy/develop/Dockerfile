FROM golang:latest

EXPOSE 7077

LABEL maintainer="GeniusRabbit Dmitry Ponomarev <demdxx@gmail.com>"
LABEL service.name=rotator
LABEL service.weight=1
LABEL service.port=8080
LABEL service.public="true"

# Connect to localhost
# $ HOSTIP=`ip -4 addr show scope global dev eth0 | grep inet | awk '{print \$2}' | cut -d / -f 1`
# $ docker run --add-host=docker:${HOSTIP} --rm -it debian

# DummyStatsD https://github.com/reiver/dummy-statsd

RUN apt-get update && apt-get install -y bzr && apt-get clean
RUN go get -u github.com/valyala/quicktemplate/qtc

# ENV PATH="$PATH:/project/bin"
# ENV GO111MODULE=on
ENV REGISTRY_DNS=http://regestry:8500/dc1?refresh_interval=5
ENV CHECK_HTTP=http://{{address}}/v1/check
ENV CHECK_INTERVAL=5s
ENV CHECK_TIMEOUT=2s

ENV AD_TRACKER_HOST=localhost:9112

# Config environment
ENV SERVICE_DOMAIN=localhost:7077
ENV STATSDCONN=metrics:8125

ENV BIGBROTHER_HOST=bigbro.service.consul:9301
ENV BIGBROTHER_PROTO=fastrpc
ENV BIGBROTHER_TIMEOUT=20

ENV EVENTSTREAM_EVENT_QUEUE_CONNECTION=nats://nats.service.consul:4222
ENV EVENTSTREAM_EVENT_QUEUE_TOPICS=events
ENV EVENTSTREAM_USERINFO_QUEUE_CONNECTION=nats://nats.service.consul:4222
ENV EVENTSTREAM_USERINFO_QUEUE_TOPICS=user_info
ENV EVENTSTREAM_WINS_QUEUE_CONNECTION=nats://nats.service.consul:4222
ENV EVENTSTREAM_WINS_QUEUE_TOPICS=wins
ENV EVENTSTREAM_READER_TOPICS=dev_adevent
ENV EVENTSTREAM_READER_FREQUENCY=100
ENV EVENTSTREAM_WRITER_TOPICS=dev_adevent
ENV EVENTSTREAM_WRITER_FREQUENCY=100

ENV ADDELTA=postgres://grbt:iFaemi9F@grdatabase.service.consul:5432/project?sslmode=disable
ENV ADSTORAGE=/var/data/grdata

# ENV BALANCE_LOCAL=redis://redis.service.consul:6379/14?idle=5&maxcon=50&timeout=3&wait=0
# ENV BALANCE_REMOTE=redis://redis.service.consul:6379/15?idle=5&maxcon=50&timeout=3&wait=0
# ENV DBCONN_BASE=postgres://grbt:iFaemi9F@grdatabase.service.consul:5432/base?sslmode=disable
# ENV DBCONN_PROJECT=postgres://grbt:iFaemi9F@grdatabase.service.consul:5432/project?sslmode=disable
# ENV ADSTORAGE_CONNECTION=postgres://grbt:iFaemi9F@grdatabase.service.consul:5432/project?sslmode=disable

ENV SSP_REQUEST_MAX_PARALLEL_REQUESTS=10
ENV SSP_REQUEST_TIMEOUT=2000

## Opentracing https://www.jaegertracing.io/docs/1.10/getting-started/
# docker run --rm -it --label=service.name=jaeger -p 6831:6831/udp -p 16686:16686 jaegertracing/all-in-one:latest
ENV JAEGER_SERVICE_NAME=rotator
ENV JAEGER_DISABLED=false
ENV JAEGER_RPC_METRICS=true
ENV JAEGER_TAGS=debug=true
ENV JAEGER_SAMPLER_TYPE=probabilistic
# ENV JAEGER_SAMPLER_TYPE=remote
ENV JAEGER_SAMPLER_PARAM=1
ENV JAEGER_SAMPLER_MANAGER_HOST_PORT=
ENV JAEGER_SAMPLER_MAX_OPERATIONS=100
ENV JAEGER_SAMPLER_REFRESH_INTERVAL=1s
ENV JAEGER_REPORTER_MAX_QUEUE_SIZE=1000
ENV JAEGER_REPORTER_FLUSH_INTERVAL=100ms
ENV JAEGER_REPORTER_LOG_SPANS=true
ENV JAEGER_ENDPOINT=
ENV JAEGER_USER=
ENV JAEGER_PASSWORD=
ENV JAEGER_AGENT_HOST=jaeger.service.consul
ENV JAEGER_AGENT_PORT=6831

WORKDIR /project/geniusrabbit/rotator
