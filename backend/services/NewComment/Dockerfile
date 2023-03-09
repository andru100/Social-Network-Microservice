FROM golang:1.18-alpine

# set a directory for the app
WORKDIR /usr/src/app

# copy all the files to the container
COPY . .

RUN go mod download

#build and name
RUN go build -o /newcomment-service

EXPOSE 4005

#run executable when container is run
CMD ["/newcomment-service"]