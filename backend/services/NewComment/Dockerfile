FROM golang:1.18-alpine

ARG HOSTIP

WORKDIR /usr/src/app

COPY . .

RUN HOSTIP=${HOSTIP}

RUN go mod download

RUN go build -o /newcomment-service

EXPOSE 4005

CMD ["/newcomment-service"]