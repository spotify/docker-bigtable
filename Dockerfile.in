## -*- docker-image-name: "spotify/bigtable-emulator" -*-

FROM golang:1.7
MAINTAINER Robert Gruener "robertg@spotify.com"

# Install netcat
RUN apt-get update
RUN apt-get install -y netcat

# Get bigtable go package
RUN go get -u cloud.google.com/go/bigtable

ADD bigtable-server.go /go/bin/bigtable-server.go
RUN go build /go/bin/bigtable-server.go
ENTRYPOINT ["/go/bigtable-server"]
EXPOSE 8080
