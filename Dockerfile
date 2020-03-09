FROM golang:1.13.8-alpine3.11 AS build
MAINTAINER renzo "luyang.sun@outlook.com"

WORKDIR $GOPATH/src/github.com/louisun/vinki
COPY . .
ENV GOPROXY=https://goproxy.io
RUN go build -o vinki .


FROM alpine:latest AS prod
COPY --from=build /go/src/github.com/louisun/vinki /vinki
WORKDIR /vinki
EXPOSE 8080
ENTRYPOINT ["./vinki"]
