package configure

import (
	"strings"
	"yfiddler/hooks"
)

type Request struct {
	HostType   int                 `yaml:"host_type"`
	Hosts      []string            `yaml:"hosts,omitempty"`
	FilterURI  []string            `yaml:"filters,omitempty"`
	FilterType int                 `yaml:"filter_type"`
	Headers    map[string][]string `yaml:"headers,omitempty"`
	DataType   int                 `yaml:"datatype"`
	Data       string              `yaml:"data,omitempty"`
	DataHooks  []hooks.Hook        `yaml:"data_hooks,omitempty"`
	Response   *Response           `yaml:"response,omitempty"`
}

func (self *Request) InHosts(host string) bool {
	if len(self.Hosts) == 0 { // 未指定 hosts  == true
		return true
	}
	for _, k := range self.Hosts {
		if k == host {
			return true
		}
	}
	return false
}

func (self *Request) Filter(filter string) bool {
	flen := len(self.FilterURI)
	if flen == 0 { //无filter 过滤要求
		return true
	}

	// 计算匹配成功次数
	n := 0
	for _, k := range self.FilterURI {
		if strings.Index(filter, k) > -1 {
			n = n + 1
		}
	}
	switch self.FilterType {
	case FILTER_TYPE_OR: // 或类型 n =0 一次没有匹配成功 return false
		if n == 0 {
			return false
		}
	case FILTER_TYPE_AND: // 与类型，n != flen  return false
		if n != flen {
			return false
		}
	case FILTER_TYPE_NO: // 非类型 , n > 0 return false
		if n > 0 {
			return false
		}
	}

	return true
}
