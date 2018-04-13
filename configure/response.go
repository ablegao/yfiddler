package configure

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func NewResponse(r *http.Request) (res *http.Response) {
	res = &http.Response{}
	res.Request = r
	res.TransferEncoding = r.TransferEncoding
	return
}

func NewByteResponse(r *http.Request, contentType string, status int, body []byte) *http.Response {
	resp := NewResponse(r)
	resp.Header = make(http.Header)
	resp.Header.Add("Content-Type", contentType)
	resp.StatusCode = status
	buf := bytes.NewBuffer(body)
	resp.ContentLength = int64(buf.Len())
	resp.Body = ioutil.NopCloser(buf)
	return resp
}

type Response struct {
	HostType   int                 `yaml:"host_type"`
	Hosts      []string            `yaml:"hosts,omitempty"`
	FilterURI  []string            `yaml:"filters,omitempty"`
	FilterType int                 `yaml:"filter_type"`
	Headers    map[string][]string `yaml:"headers,omitempty"`
	DataType   int                 `yaml:"datatype,omitempty"`
	Data       string              `yaml:"data,omitempty"`
}

func (self *Response) GenReadCloser() (int64, io.ReadCloser, error) {
	switch self.DataType {
	case DATA_TYPE_TEXT:
		buf := bytes.NewBufferString(self.Data)
		return int64(buf.Len()), ioutil.NopCloser(buf), nil
	case DATA_TYPE_FILE:
		f, err := os.Open(self.Data)
		if err != nil {
			return 0, nil, err
		}
		fi, err := f.Stat()
		if err != nil {
			return 0, nil, err
		}
		return fi.Size(), io.ReadCloser(f), nil

	}

	return 0, nil, errors.New("Gen ERROR NULL")
}

func (self *Response) InHosts(host string) bool {
	n := 0
	for _, k := range self.Hosts {
		if k == host {
			n = n + 1
		}
	}
	switch self.HostType {
	case HOST_TYPE_OR:
		if n == 0 {
			return false
		}
	}
	return true
}

func (self *Response) Filter(filter string) bool {
	n := 0
	for _, k := range self.FilterURI {
		if strings.Index(filter, k) > -1 {
			n = n + 1
		}
	}

	switch self.FilterType {
	case FILTER_TYPE_NONE:
		return true
	case FILTER_TYPE_OR:
		if n == 0 {
			return false
		}
	case FILTER_TYPE_AND:
		if n != len(self.FilterURI) {
			return false
		}
	case FILTER_TYPE_NO:
		if n > 0 {
			return false
		}
	}

	return true
}
