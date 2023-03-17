FROM golang:1.18-alpine

ARG SENDER_EMAIL

WORKDIR /usr/src/app

COPY . .

ENV SENDER_EMAIL=${SENDER_EMAIL}

RUN go mod download

RUN go build -o /requestotp-service

EXPOSE 4011

CMD ["/requestotp-service"]