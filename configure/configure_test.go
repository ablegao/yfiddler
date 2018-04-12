package configure

import (
	"io/ioutil"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

func Test_loadyaml(t *testing.T) {
	yf, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(yf))
	conf := new(Configure)
	err = yaml.Unmarshal(yf, conf)
	if err != nil {
		t.Error(err)
	}

	t.Log(conf)

}
