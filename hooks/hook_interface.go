package hooks

import log "github.com/sirupsen/logrus"

type Hook struct {
	Name   string   `yaml:"name,omitempty"`
	Args   []string `yaml:"args,omitempty"`
	Plugin string   `yaml:"plugin,omitempty"`
}

func (self Hook) Gen(body string) string {
	p, ok := pluginsMap[self.Name]
	if !ok {
		log.Error("Hook plugin: ", self.Name, " not exists!")
		return body
	}
	args := []string{body}
	args = append(args, self.Args...)
	body = p.Gen(args...)
	return body
}

type HookPlugin interface {

	// Gen 中， 第一个参数将作为原始数据处理，yaml 中Hook的其他参数按顺序追加
	Gen(args ...string) string
}
