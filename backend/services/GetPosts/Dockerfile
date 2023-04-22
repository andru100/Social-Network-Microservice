FROM golang:1.18-alpine

# set a directory for the app
WORKDIR /usr/src/app


# copy all the files to the container
COPY . .

RUN go mod download

#build and name
RUN go build -o /getposts-service

EXPOSE 4009

#run executable when container is run
CMD ["/getposts-service"]