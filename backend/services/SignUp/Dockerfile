FROM golang:1.18-alpine

ARG HOSTIP
ARG LOCK_DURATION_1
ARG LOCK_DURATION_2
ARG LOCK_DURATION_3

# set a directory for the app
WORKDIR /usr/src/app

ENV HOSTIP=${HOSTIP}
ENV LOCK_DURATION_1=${LOCK_DURATION_1}
ENV LOCK_DURATION_2=${LOCK_DURATION_2}
ENV LOCK_DURATION_3=${LOCK_DURATION_3}

# copy all the files to the container
COPY . .

RUN go mod download

#build and name
RUN go build -o /signup-service

EXPOSE 4002

#run executable when container is run
CMD ["/signup-service"]