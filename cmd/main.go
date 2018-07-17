package main

import (
	"flag"
	"os"
	"simple_collector"
)

var (
	sourceFile string
	outputFile string
	outputFileType string
)

func init() {
	flag.StringVar(&sourceFile, "source", "source.txt", "path to url sourceFile file")
	flag.StringVar(&outputFile, "output", "output.txt", "path output file")
	flag.StringVar(&outputFileType, "t", "txt", "type of output file")
	flag.Parse()

	Environment()
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
}

func main() {
	simple_collector.Run(sourceFile, outputFile, outputFileType)
}