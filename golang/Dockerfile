FROM golang:latest as DEV
RUN mkdir /go/src/app
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOCS=linux GOACH=amd64 go build main.go
RUN chmod a+x main

FROM golang:alpine as PROD
WORKDIR /root/
COPY --from=DEV /go/src/app/main .
CMD ["./main"]
