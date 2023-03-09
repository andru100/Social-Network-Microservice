FROM golang:1.18-alpine

WORKDIR /usr/src/app

COPY . .

RUN go mod download

RUN go build -o /backend-graphql

EXPOSE 8080 

CMD ["/backend-graphql"]