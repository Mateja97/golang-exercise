FROM golang:1.23.4 AS build


WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app

EXPOSE 8080 8080
ENTRYPOINT ["/build/app"]
