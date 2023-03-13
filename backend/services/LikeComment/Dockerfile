FROM golang:1.18-alpine

ARG HOSTIP

WORKDIR /usr/src/app

COPY . .

RUN HOSTIP=${HOSTIP}

RUN go mod download


RUN go build -o /likecomment-service

EXPOSE 4003


CMD ["/likecomment-service"]