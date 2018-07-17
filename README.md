# Simple Collector

Don't forget to add `$GOPATH/bin` to your `$PATH`.
### Step 1: clone simple_collector
### Step 2: install dep (if not installed)
```bash
go get -u github.com/golang/dep/cmd/dep
```

### Running
```bash
go build && ./cmd -source source.txt -output a.json -t json
```