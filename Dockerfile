FROM golang:1.8.7-stretch

WORKDIR /go/src/github.com/caiguanhao/octopus

ADD . .

ENV LIBRARY_PATH /go/src/github.com/caiguanhao/octopus

RUN make package
