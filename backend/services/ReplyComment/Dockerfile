FROM golang:1.18-alpine

ARG HOSTIP

WORKDIR /usr/src/app

COPY . .

ENV HOSTIP=${HOSTIP}

RUN go mod download


RUN go build -o /replycomment-service

EXPOSE 4004

CMD ["/replycomment-service"]