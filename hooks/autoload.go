package hooks

import (
	"fmt"
	"io/ioutil"
	"plugin"
	"strings"

	log "github.com/sirupsen/logrus"
)

var pluginsMap = map[string]HookPlugin{}

func GetHook(name string) HookPlugin {
	if fun, ok := pluginsMap[name]; ok {
		return fun
	} else {
		return nil
	}
}

func Autoload(pdir string) {
	files, err := ioutil.ReadDir(pdir)
	if err != nil {
		log.Error(err)
		return
	}
	for _, file := range files {
		fileInfo := strings.Split(file.Name(), ".")
		if !file.IsDir() && len(fileInfo) == 2 && fileInfo[1] == "so" {
			log.Debug("Plugin file: ", file.Name())
			p, err := plugin.Open(fmt.Sprintf("%s/%s", pdir, file.Name()))
			if err != nil {
				log.Error(err)
				continue
			}
			//	pluginsMap[fileInfo[0]] =
			plug, err := p.Lookup(strings.Title(fileInfo[0]))
			if err != nil {
				log.Error(err)
				continue
			}
			pluginsMap[fileInfo[0]] = plug.(HookPlugin)
		}
	}
	log.Debug(pluginsMap)
}
