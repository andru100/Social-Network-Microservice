FROM golang:1.18-alpine

# set a directory for the app
WORKDIR /usr/src/app

# copy all the files to the container
COPY . .

RUN go mod download

#build and name
RUN go build -o /likecomment-service

EXPOSE 4003

#run executable when container is run
CMD ["/likecomment-service"]