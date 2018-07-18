package main

import (
	"flag"
	"os"
	"simple_collector"
	"github.com/vidmed/logger"
	"strconv"
)

var (
	sourceFile string
	outputFile string
	outputFileType string
	logLevel int
)

func init() {
	flag.StringVar(&sourceFile, "source", "source.txt", "path to url sourceFile file")
	flag.StringVar(&outputFile, "output", "output.txt", "path output file")
	flag.StringVar(&outputFileType, "t", "txt", "type of output file")
	flag.IntVar(&logLevel, "loglvl", 2, "log level panic = 0, fatal = 1, error = 2, warning = 3, info = 4, debug = 5")
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
}

func main() {
	logger.Get().Infof("Application running with parameters: source = %s, output = %s, t = %s",
		sourceFile, outputFile, outputFileType)
	simple_collector.Run(sourceFile, outputFile, outputFileType)
}