#Compile stage
FROM golang:1.11.4-alpine3.8 AS build-env
ENV CGO_ENABLED 0
RUN apk add --no-cache git
ADD . /go/src/notification-ms

RUN go get -u github.com/revel/revel
RUN go get -u github.com/revel/cmd/revel
RUN go get -u github.com/go-redis/redis
RUN go get -u github.com/prometheus/common/log

#build revel app
RUN revel build notification-ms app dev

# Final stage
FROM alpine:3.8
EXPOSE 9000
WORKDIR /
COPY --from=build-env /go/app /
ENTRYPOINT /run.sh