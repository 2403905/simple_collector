package simple_collector

import (
	"fmt"
	"net/http"
	"time"
	"golang.org/x/net/context/ctxhttp"
	"context"
)

var (
	maxWorkerCount = int(100)
)

func Run(sourceFile, outputFile, outputFileType string) {
	jobs := make(chan string, maxWorkerCount)
	results := make(chan *ResponseData, maxWorkerCount)

	output, err := NewOutput(getSaver(outputFileType))
	if err != nil {
		fmt.Println(err)
		return
	}

	u, err := getJobs(sourceFile)
	if err != nil {
		fmt.Println(err)
		return
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
	err = output.saveResult(outputFile)
	if err != nil {
		fmt.Println(err)
	}
}

func worker(jobs <-chan string, results chan<- *ResponseData) {
	for j := range jobs {
		res := &ResponseData{Url:j}
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

func getSaver(t string) (func(resp []*ResponseData, outputFile string) error) {
	switch t {
	case "txt":
		return saveText
	case "json":
		return saveJson
	default:
		fmt.Printf("ouput type %s is not supported \n", t)
		return nil
	}
}
