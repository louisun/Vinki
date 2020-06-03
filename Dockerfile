FROM mhart/alpine-node:14 AS front-build
WORKDIR /frontend
COPY ./frontend .
RUN yarn config set registry https://registry.npm.taobao.org
RUN yarn install
RUN yarn run build

FROM golang:1.14-alpine AS build
MAINTAINER renzo "luyang.sun@outlook.com"

WORKDIR $GOPATH/src/github.com/louisun/vinki
COPY . .
COPY --from=front-build /frontend/build ./statics
ENV GOPROXY=https://goproxy.io
RUN go get github.com/rakyll/statik
RUN statik -src=statics -include="*.html,*.js,*.json,*.css,*.png,*.svg,*.ico,*.woff,*.woff2,*.txt" -f
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --update gcc g++
RUN go build -o vinki .

FROM alpine:latest AS prod
COPY --from=build /go/src/github.com/louisun/vinki/vinki /vinki/
WORKDIR /vinki
EXPOSE 6166
ENTRYPOINT ["./vinki"]
