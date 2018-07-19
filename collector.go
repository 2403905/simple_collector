package simple_collector

import (
	"net/http"
	"time"
	"golang.org/x/net/context/ctxhttp"
	"context"
	"github.com/vidmed/logger"
	"os"
)

var (
	maxWorkerCount = int(10000)
)

func Run(sourceFile, outputFile, outputFileType string) {
	jobs := make(chan string, maxWorkerCount)
	results := make(chan *ResponseData, maxWorkerCount)

	// init output method
	output, err := NewOutput(getWriter(outputFileType))
	if err != nil {
		logger.Get().Errorln(err)
		os.Exit(1)
	}

	u, err := getJobs(sourceFile)
	if err != nil {
		logger.Get().Errorln(err)
		os.Exit(1)
	}

	if maxWorkerCount > len(u) {
		maxWorkerCount = len(u)
	}

	for w := 1; w <= maxWorkerCount; w++ {
		go worker(jobs, results)
	}

	for _, j := range u {
		jobs <- j
	}
	close(jobs)

	data := make([]*ResponseData, len(u))

	for k := range data {
		r := <-results
		data[k] = r
	}

	output.ResponseData = data
	err = output.writeResult(outputFile)
	if err != nil {
		logger.Get().Errorln(err)
		os.Exit(1)
	}
	logger.Get().Info("Well done.")
}

func worker(jobs <-chan string, results chan<- *ResponseData) {
	for j := range jobs {
		res := &ResponseData{Url: j}
		resp, latency, err := Send(j)
		if err != nil {
			res.Error = err.Error()
		} else {
			res.Headers = resp.Header
			res.ResponseCode = resp.StatusCode
			res.Latency = latency.String()
		}

		results <- res
	}
}

func Send(url string) (response *http.Response, latency time.Duration, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := http.DefaultClient

	start := time.Now()
	response, err = ctxhttp.Get(ctx, client, url)
	latency = time.Since(start)

	return
}
