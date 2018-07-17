package simple_collector

import (
	"os"
	"bufio"
	"github.com/pkg/errors"
	"fmt"
)

func getJobs(sourceFile string) ([]string, error) {
	file, err := os.Open(sourceFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		return nil, errors.Wrap(scanner.Err(), fmt.Sprintf("error during parsing file %s ", sourceFile))
	}
	return lines, nil
}