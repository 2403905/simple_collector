FROM golang:1.10.3

WORKDIR /go/src/simple_collector
COPY Gopkg.lock Gopkg.toml /go/src/simple_collector/
RUN go get -u github.com/golang/dep/...
RUN dep ensure -vendor-only

COPY . /go/src/simple_collector/

WORKDIR cmd
RUN go build -o  /go/bin/simple_collector

CMD ["/go/bin/simple_collector"]
