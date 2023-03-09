FROM golang:1.18-alpine

# set a directory for the app
WORKDIR /usr/src/app

# copy all the files to the container
COPY . .

RUN go mod download

#build and name
RUN go build -o /replycomment-service

EXPOSE 4004

#run executable when container is run
CMD ["/replycomment-service"]