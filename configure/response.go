package configure

import (
	"bytes"
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
	HostType   int               `yaml:"host_type"`
	Hosts      []string          `yaml:"hosts"`
	FilterURI  []string          `yaml:"filters"`
	FilterType int               `yaml:"filter_type"`
	Headers    map[string]string `yaml:"headers"`
	DataType   int               `yaml:"datatype"`
	Data       string            `yaml:"data"`
}

func (self *Response) GetResponse(r *http.Request, res *http.Response) bool {
	res.Request = r
	res.TransferEncoding = r.TransferEncoding
	res.Header = make(http.Header)
	for k, v := range self.Headers {
		res.Header.Set(k, v)
	}
	switch self.DataType {
	case DATA_TYPE_TEXT:
		buf := bytes.NewBufferString(self.Data)
		res.ContentLength = int64(buf.Len())
		res.Body = ioutil.NopCloser(buf)
	case DATA_TYPE_FILE:
		f, err := os.Open(self.Data)
		if err != nil {
			panic(err)
			return false
		}
		if fi, err := f.Stat(); err == nil {
			res.ContentLength = fi.Size()
			res.Body = f
		} else {
			panic(err)
			return false
		}
	}

	return true
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
