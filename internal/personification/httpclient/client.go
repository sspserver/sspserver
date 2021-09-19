package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/sspserver/udetect"
)

const (
	macroIP        = "{ip}"
	macroUserAgent = "{ua}"
)

// Client for user specific detection
type Client struct {
	url             string
	method          string
	isPreparableURL bool
	client          *http.Client
}

// NewClient returns client with HTTP protocol support
func NewClient(url string, options ...Option) *Client {
	cli := &Client{
		url:             url,
		method:          http.MethodPost,
		isPreparableURL: strings.Contains(url, macroIP) || strings.Contains(url, macroUserAgent),
		client:          &http.Client{},
	}
	for _, opt := range options {
		opt(cli)
	}
	return cli
}

// Detect information from request
func (cli *Client) Detect(req *udetect.Request) (resp *udetect.Response, err error) {
	var (
		httpRequest  *http.Request
		httpResponse *http.Response
		body         io.Reader
	)
	if cli.method == http.MethodPost {
		var data bytes.Buffer
		json.NewEncoder(&data).Encode(req)
		body = &data
	}
	httpRequest, err = http.NewRequest(cli.method, cli.preparedURL(req), body)
	if err != nil {
		return nil, err
	}
	httpResponse, err = cli.client.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	resp = &udetect.Response{}
	err = json.NewDecoder(httpResponse.Body).Decode(resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (cli *Client) preparedURL(req *udetect.Request) string {
	if !cli.isPreparableURL {
		return cli.url
	}
	replacer := strings.NewReplacer(macroIP, req.Ip, macroUserAgent, req.Ua)
	return replacer.Replace(cli.url)
}
