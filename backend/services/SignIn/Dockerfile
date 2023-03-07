FROM golang:1.18-alpine

# set a directory for the app
WORKDIR /usr/src/app

COPY go.mod .
COPY go.sum .
RUN go mod download

# copy all the files to the container
COPY . .

#build and name
RUN go build -o '/signin-service'

EXPOSE 4005

#run executable when container is run
CMD ["/signin-service"]