FROM golang:1.12 AS build-env
RUN mkdir /autodeploy
WORKDIR /autodeploy
COPY go.mod .
COPY go.sum .

RUN CGO_ENABLED=0 go mod download
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/auto

FROM alpine
RUN mkdir -p /go/bin/config
WORKDIR /go/bin
COPY ./config/docker.json /go/bin/config/conf.json
COPY --from=build-env /go/bin/auto /go/bin/auto
ENTRYPOINT ["/go/bin/auto"]
