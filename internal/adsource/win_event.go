//
// @project GeniusRabbit rotator 2016 – 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2017
//

package adsource

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pierrec/lz4"
)

// Errors
var (
	ErrInvalidURL = errors.New("Invalid URL value")
)

// WinEvent models
type WinEvent struct {
	Time     time.Time   `json:"t,omitempty"`
	Counter  int         `json:"c,omitempty"`
	URL      string      `json:"u"`
	Data     interface{} `json:"d,omitempty"`
	Response string      `json:"-"`
	Error    error       `json:"-"`
	Status   int         `json:"-"`
}

// WinEventDecode from bytes
func WinEventDecode(data string, ext ...interface{}) (*WinEvent, error) {
	var win = &WinEvent{}
	if len(ext) > 0 {
		win.Data = ext[0]
	}
	return win, win.Decode(data)
}

// Validate event object
func (e *WinEvent) Validate() error {
	url, err := url.Parse(e.URL)
	if err != nil {
		return err
	}
	if len(url.Scheme) < 1 || len(url.Host) < 5 {
		return ErrInvalidURL
	}
	return nil
}

// Inc event counter
func (e *WinEvent) Inc() *WinEvent {
	e.Time = time.Now()
	e.Counter++
	return e
}

// Reset event counter
func (e *WinEvent) Reset() *WinEvent {
	e.Time = time.Now()
	e.Counter = 0
	return e
}

// Ping client server request about win
func (e *WinEvent) Ping() bool {
	if strings.HasPrefix(e.URL, "//") {
		e.URL = "http:" + e.URL
	}

	var resp *http.Response
	if u, err := url.QueryUnescape(e.URL); nil == err {
		e.URL = u
	}

	resp, e.Error = http.Get(e.URL)
	if resp != nil && e.Error == nil {
		e.Status = resp.StatusCode
		if data, _ := ioutil.ReadAll(resp.Body); nil != data {
			e.Response = string(data)
		} else {
			e.Response = ""
		}
		resp.Body.Close()
	}
	return resp != nil && e.Error == nil && http.StatusOK == resp.StatusCode
}

// CanContinue try
func (e *WinEvent) CanContinue() bool {
	return (e.Status >= http.StatusOK && e.Status < http.StatusBadRequest) || e.Status >= http.StatusInternalServerError
}

// Encode URL object
func (e *WinEvent) Encode() (string, error) {
	var (
		buff   bytes.Buffer
		writer = lz4.NewWriter(&buff)
		err    = json.NewEncoder(writer).Encode(e)
	)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buff.Bytes()), nil
}

// Decode URL object
func (e *WinEvent) Decode(data string) error {
	var decoded, err = base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}

	return json.NewDecoder(
		lz4.NewReader(bytes.NewReader(decoded)),
	).Decode(e)
}
