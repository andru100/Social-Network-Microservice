FROM golang:1.18-alpine

ARG SENDER_EMAIL
ARG HOSTIP
ENV HOSTIP=${HOSTIP}

WORKDIR /usr/src/app

COPY . .

ENV SENDER_EMAIL=${SENDER_EMAIL}

ENV HOSTIP=${HOSTIP}

RUN go mod download

RUN go build -o /secureupdate-service

EXPOSE 4012

CMD ["/secureupdate-service"]