package simple_collector

import (
	"net/http"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"github.com/vidmed/logger"
)

type ResponseData struct {
	Url          string
	ResponseCode int         `json:"response_code"`
	Latency      string      `json:"latency"`
	Headers      http.Header `json:"headers"`
	Error        string      `json:"error"`
}

type Output struct {
	ResponseData []*ResponseData
	saver        func(resp []*ResponseData, outputFile string) error
}

func NewOutput(saver func(resp []*ResponseData, outputFile string) error) (*Output, error)  {
	if saver == nil {
		return nil, fmt.Errorf("saver not set")
	}
	return &Output{saver:saver}, nil
}

func (s Output) saveResult(outputFile string) error {
	if len(s.ResponseData) == 0 {
		return fmt.Errorf("there nothing to save")
	}

	return s.saver(s.ResponseData, outputFile)
}

func saveJson(resp []*ResponseData, outputFile string) (err error) {
	data, err := json.Marshal(resp)
	if err != nil {
		logger.Get().Errorln(err)
		return
	}
	err = ioutil.WriteFile(outputFile, data, 0644)
	if err != nil {
		logger.Get().Errorln(err)
		return
	}
	return
}

func saveText(resp []*ResponseData, outputFile string) (err error) {
	f, err := os.Create(outputFile)
	if err != nil {
		logger.Get().Errorln(err)
		return
	}
	defer f.Close()

	var str string
 	for _, r :=range resp {
		str = str + fmt.Sprintf(
			"url:%s code:%d, latency%s, headers:%v, error:%s \n",
			r.Url, r.ResponseCode, r.Latency, r.Headers, r.Error)
	}

	_, err =f.Write([]byte(str))
	if err != nil {
		logger.Get().Errorln(err)
		return
	}
	return
}
