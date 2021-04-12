FROM golang:1.16 AS build
WORKDIR /go/src
COPY . .

ENV CGO_ENABLED=0
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o webdav-submission .

FROM scratch AS runtime
ENV GIN_MODE=release
COPY --from=build /go/src/webdav-submission ./
ENTRYPOINT ["./webdav-submission"]
