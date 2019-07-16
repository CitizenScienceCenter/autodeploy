FROM golang:1.12-alpine
RUN mkdir /auto
WORKDIR /auto
COPY go.mod . # <- COPY go.mod and go.sum files to the workspace
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/auto 

FROM scratch
COPY --from=build-env /go/bin/auto /go/bin/auto
ENTRYPOINT ["/go/bin/auto"]
