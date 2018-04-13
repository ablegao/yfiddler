package configure

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

const (
	DATA_TYPE_TEXT = 0
	DATA_TYPE_FILE = 1
	DATA_TYPE_NONE = 2

	FILTER_TYPE_OR   = 0 //或
	FILTER_TYPE_AND  = 1 //和
	FILTER_TYPE_NO   = 2 //非
	FILTER_TYPE_NONE = 3

	HOST_TYPE_OR   = 0
	HOST_TYPE_NONE = 1
)

type Configure struct {
	Requests  []*Request  `yaml:"requests,omitempty"`
	Responses []*Response `yaml:"responses,omitempty"`
}

func NewConfigureByYaml(f string) (*Configure, error) {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	conf := new(Configure)
	err = yaml.Unmarshal(b, conf)
	return conf, err
}
