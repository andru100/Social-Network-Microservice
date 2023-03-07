FROM golang:1.18-alpine

# set a directory for the app
WORKDIR /usr/src/app

COPY go.mod .
COPY go.sum .
RUN go mod download

# copy all the files to the container
COPY . .

#build and name
RUN go build -o '/postfile-service'

EXPOSE 4002

#run executable when container is run
CMD ["/postfile-service"]