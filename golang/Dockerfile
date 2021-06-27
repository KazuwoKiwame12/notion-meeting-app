FROM golang:1.16 as DEV

ENV GOCS=linux
ENV GOACH=amd64
ENV CGO_ENABLED=0
ENV GO111MODULE=on

WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build server.go

FROM alpine:3.14 as PROD
COPY --from=DEV /go/src/app/server /go/src/app/asset/slack/modalview.json ./
ENTRYPOINT ["./server"]