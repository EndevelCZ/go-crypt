# https://medium.com/@petomalina/using-go-mod-download-to-speed-up-golang-docker-builds-707591336888
FROM golang:alpine AS build-env
RUN apk add bash ca-certificates git gcc g++ libc-dev 
ENV GO111MODULE=on
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/hello
FROM alpine AS last-phase
LABEL maintainer="Adam Plánský <adamplansky@gmail.com>"
RUN apk add ca-certificates
# Finally we copy the statically compiled Go binary.
COPY --from=build-env /go/bin/hello /bin/hello
CMD ["/bin/hello"]

