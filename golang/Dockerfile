FROM golang:latest

RUN mkdir /go/src/app
WORKDIR /go/src/app
COPY . .
RUN go mod download

EXPOSE 8080