package utils

import (
	"destroyer-monitor/config"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type HttpConPool struct {
	*config.HttpCli
	cli *http.Client
}

func NewHttpCli(conf *config.HttpCli) *HttpConPool {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        conf.MaxIdleConns,
			MaxIdleConnsPerHost: 50,
		},
		Timeout: time.Duration(conf.Timeout) * time.Millisecond,
	}
	return &HttpConPool{
		conf,
		client,
	}
}

func (h HttpConPool) ReqGetDiscallResponse(url string) {
	_, _ = h.cli.Get(url)
}

func (h HttpConPool) ReqGet(reqUrl string) (string, int, error) {
	u, err := url.Parse(reqUrl)
	if nil != err {
		return "http client request url is illegal.", 500, err
	}
	u.RawQuery = u.Query().Encode()

	response, err := h.cli.Get(reqUrl)
	if nil != err {
		return "http client request do error.", 500, err
	} else if nil != response {
		body := []byte{}
		if response.ContentLength < 16384 {
			body, _ = ioutil.ReadAll(response.Body)
			response.Body.Close()
		} else {
			body = []byte("content greater than 16kb")
			fmt.Println("[", time.Now().Format("2006-01-02 15:04:05"), "] ReqGet to advertiser response.")
		}
		return ByteToStr(body), response.StatusCode, nil
	} else {
		return "http client request do res nil.", 500, errors.New("http client request do res nil.")
	}
}


