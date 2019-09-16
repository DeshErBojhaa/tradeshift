
# #build stage
# FROM golang:alpine AS builder
# WORKDIR /go/src/app
# COPY . .
# RUN apk add --no-cache git
# RUN go get -d -v ./...
# RUN go install -v ./...

#final stage
FROM golang:alpine
COPY ./appbinary ./
ENTRYPOINT ./appbinary
LABEL Name=tradeshift Version=0.0.1
EXPOSE 8080
