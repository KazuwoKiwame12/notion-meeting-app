FROM golang:latest as DEV
RUN mkdir /go/src/app
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOCS=linux GOACH=amd64 go build server.go

FROM golang:alpine as PROD
WORKDIR /root/
COPY --from=DEV /go/src/app/server .
ENTRYPOINT ["./server"]
# WORKDIR /go/src/app
# COPY --from=DEV /go/src/app .
# ENTRYPOINT ["go", "run", "server.go"]
#