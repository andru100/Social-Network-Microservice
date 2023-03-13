FROM golang:1.18-alpine

ARG HOSTIP

WORKDIR /usr/src/app

COPY . .

RUN HOSTIP=${HOSTIP}

RUN go mod download

RUN go build -o /backend-graphql

EXPOSE 8080 

CMD ["/backend-graphql"]