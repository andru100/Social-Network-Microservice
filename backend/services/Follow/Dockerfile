FROM golang:1.18-alpine

ARG HOSTIP

WORKDIR /usr/src/app

COPY . .

ENV HOSTIP=${HOSTIP}

RUN go mod download


RUN go build -o /follow-service

EXPOSE 4013


CMD ["/follow-service"]