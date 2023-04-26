FROM golang:1.18-alpine

ARG SENDER_EMAIL

ARG HOSTIP

ARG ENABLE_AWS

WORKDIR /usr/src/app

COPY . .

ENV SENDER_EMAIL=${SENDER_EMAIL}

ENV HOSTIP=${HOSTIP}

ENV ENABLE_AWS=${ENABLE_AWS}

RUN go mod download

RUN go build -o /requestotp-service

EXPOSE 4011

CMD ["/requestotp-service"]