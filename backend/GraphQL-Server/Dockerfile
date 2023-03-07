FROM golang:1.18-alpine

WORKDIR /usr/src/app

COPY . .

RUN go mod download

RUN go build -o /Backend

EXPOSE 8080 8081

CMD ["/Backend"]