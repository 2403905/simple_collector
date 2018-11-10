package simple_collector

import (
	"net/http"
	"time"
	"github.com/vidmed/logger"
	"os"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/53.0.2785.143 " +
		"Safari/537.36"
)

var (
	requestTimeout = 5
)

func Run(sourceFile, outputFile, outputFileType string, maxWorkerCount, timeout int) {
	// init timeout
	requestTimeout = timeout

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

	jobs := make(chan string, len(u))
	results := make(chan *ResponseData, maxWorkerCount)
	countTimeoutErr := make(chan int64, maxWorkerCount)

	go func(countTimeoutErr chan int64) {
		logger.Get().Infoln("Context error counter ran")
		var count, total int64
		var ok bool
	BrakeLoop:
		for {
			select {
			case count, ok = <-countTimeoutErr:
				if !ok {
					break BrakeLoop
				}
				total += count
			default:
			}
		}

		logger.Get().Infof("context error count %d", total)
	}(countTimeoutErr)

	for w := 1; w <= maxWorkerCount; w++ {
		go worker(jobs, results, countTimeoutErr)
	}

	for _, j := range u {
		jobs <- j
	}
	close(jobs)
	logger.Get().Infof("%d jobs sent.", len(u))

	data := make([]*ResponseData, len(u))

	for k := range data {
		r := <-results
		data[k] = r
	}
	logger.Get().Infoln("Getting from result channel done")

	close(countTimeoutErr)

	logger.Get().Infoln("Response data writer ran")
	output.ResponseData = data
	err = output.writeResult(outputFile)
	if err != nil {
		logger.Get().Errorln(err)
		os.Exit(1)
	}
	logger.Get().Info("Well done.")
}

func worker(jobs <-chan string, results chan<- *ResponseData, countTimeoutErr chan<- int64) {
	for j := range jobs {
		res := &ResponseData{Url: j}
		resp, latency, err := Send(j)
		if err != nil {
			res.Error = err.Error()
			countTimeoutErr <- 1
		} else {
			res.ResponseCode = resp.StatusCode
		}
		res.Latency = latency.String()

		results <- res
	}
}

func Send(url string) (response *http.Response, latency time.Duration, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(requestTimeout)*time.Second)
	defer cancel()
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", userAgent)

	start := time.Now()
	response, err = ctxhttp.Do(ctx, client, req)
	latency = time.Since(start)

	return
}
