# Simple Collector

### Installation with docker
1. Install docker and docker-compose
2. Clone repository simple_collector
3. Go to repository directory
4. Prepare i/o directory
   ```bash
   mkdir -p /tmp/simple_collector
   cp source/source.txt /tmp/simple_collector/
   ```
5. Run
   ```bash
   docker-compose up
   ```
6. Output
   ```bash
   cat /tmp/simple_collector/output.txt
   ```


### Installation without docker
1. Install golang
2. Add `$GOPATH/bin` to your `$PATH`.
3. Clone repository simple_collector
4. Install dep (if not installed)
   ```bash
   go get -u github.com/golang/dep/cmd/dep
   ```
4. Go to repository directory
5. Run dep
   ```bash
   dep ensure -v
   ```
6. Run
   ```bash
   cd cmd
   go build && ./cmd -source source.txt -output output.json -t json
   ```
7. Output
   ```bash
   cat output.json
   ```