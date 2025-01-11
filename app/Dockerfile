FROM golang:1.23.4 AS build
WORKDIR /build

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin v1.54.2

COPY go.mod go.sum ./
COPY .golangci.yml ./
RUN go mod tidy
RUN go mod vendor

COPY . .

RUN golangci-lint run ./... || echo "Linting issues found."

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app


FROM golang:1.23.4

WORKDIR /app

COPY --from=build /build/app .

EXPOSE 8080
ENTRYPOINT ["/app/app"]
