package main

import (
	"flag"
	"os"
	"simple_collector"
	"github.com/vidmed/logger"
	"strconv"
	"time"
)

var (
	sourceFile     string
	outputFile     string
	outputFileType string
	workersCount   int
	timeout        int
	logLevel       int
)

func init() {
	flag.StringVar(&sourceFile, "source", "source.txt", "path to url sourceFile file")
	flag.StringVar(&outputFile, "output", "output.txt", "path output file")
	flag.StringVar(&outputFileType, "t", "txt", "type of output file")
	flag.IntVar(&workersCount, "wc", 1000, "workersCount")
	flag.IntVar(&timeout, "timeout", 60, "request timeout (second)")
	flag.IntVar(&logLevel, "loglvl", 4, "log level panic = 0, fatal = 1, error = 2, warning = 3, info = 4, debug = 5")
	flag.Parse()

	Environment()

	// Init Logger
	logger.Init(logLevel)
	logger.AddStackHook()
}

func Environment() {
	if s, ok := os.LookupEnv("SOURCE_FILE"); ok {
		sourceFile = s
	}
	if o, ok := os.LookupEnv("OUTPUT_FILE"); ok {
		outputFile = o
	}
	if t, ok := os.LookupEnv("OUTPUT_FILE_TYPE"); ok {
		outputFileType = t
	}
	if l, ok := os.LookupEnv("LOG_LEVEL"); ok {
		var err error
		logLevel, err = strconv.Atoi(l)
		if err != nil {
			logger.Get().Errorf("Can't convert environment variable LOG_LEVEL=%s to integer", l)
		}
	}
	if l, ok := os.LookupEnv("WORKERS_COUNT"); ok {
		var err error
		workersCount, err = strconv.Atoi(l)
		if err != nil {
			logger.Get().Errorf("Can't convert environment variable WORKERS_COUNT=%s to integer", l)
		}
	}
	if l, ok := os.LookupEnv("TIMEOUT"); ok {
		var err error
		timeout, err = strconv.Atoi(l)
		if err != nil {
			logger.Get().Errorf("Can't convert environment variable TIMEOUT=%s to integer", l)
		}
	}
}

func main() {
	start := time.Now()
	logger.Get().Infof("Application running with parameters: source = %s, output = %s, t = %s, workersCount = %d, timeout = %d",
		sourceFile, outputFile, outputFileType, workersCount, timeout)
	simple_collector.Run(sourceFile, outputFile, outputFileType, workersCount, timeout)
	elapsed := time.Since(start)
	logger.Get().Infof("It took %s", elapsed)
}
