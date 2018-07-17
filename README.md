# Simple Collector

Don't forget to add `$GOPATH/bin` to your `$PATH`.
### Clone simple_collector
### Install dep (if not installed)
```bash
go get -u github.com/golang/dep/cmd/dep
```
### Run
```bash
go build && ./cmd -source source.txt -output a.json -t json
```