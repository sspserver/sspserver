package httpclient

import (
	"net"
	"net/http"
	"time"
	"unsafe"
)

// Option of client connection
type Option func(cli *Client)

// WithGET method of request
func WithGET(cli *Client) {
	cli.method = http.MethodGet
}

// WithPOST method of request
func WithPOST(cli *Client) {
	cli.method = http.MethodPost
}

// WithTransport interface
func WithTransport(transport http.RoundTripper) Option {
	return func(cli *Client) {
		cli.client.Transport = transport
	}
}

// WithTimeout sets request timeouts
func WithTimeout(timeout, keepAlive time.Duration) Option {
	return func(cli *Client) {
		cli.client.Timeout = timeout

		if t := getTransport(cli.client.Transport); t != nil {
			t.DialContext = (&net.Dialer{
				Timeout:   timeout,
				KeepAlive: keepAlive,
				DualStack: true,
			}).DialContext
			cli.client.Transport = t
		}
	}
}

// WithMaxConn from the client
func WithMaxConn(connCount int) Option {
	return func(cli *Client) {
		t := mustGetTransport(cli.client.Transport)
		t.MaxConnsPerHost = connCount
		t.MaxIdleConns = connCount
		t.MaxIdleConnsPerHost = connCount
		cli.client.Transport = t
	}
}

func mustGetTransport(transport http.RoundTripper) *http.Transport {
	t := getTransport(transport)
	if t == nil {
		panic("unsupported transport object")
	}
	return t
}

func getTransport(transport http.RoundTripper) *http.Transport {
	if t, _ := transport.(*http.Transport); t != nil {
		dt, _ := http.DefaultTransport.(*http.Transport)
		if dt != nil && (unsafe.Pointer)(t) != (unsafe.Pointer)(dt) {
			return t
		}
	}
	return &http.Transport{}
}
