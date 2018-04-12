package header

import (
	"fmt"
	"io/ioutil"
	"plugin"
	"strings"

	"net/http"

	"github.com/elazarl/goproxy"
	log "github.com/sirupsen/logrus"
)

var pluginsMap = map[string]ProxyPlugin{}

type ProxyPlugin interface {
	OnStart()
	OnStop()
	Reset()
	InHost(string) bool
	Filter(string) bool
}

type ProxyPluginRequest interface {
	OnStart()
	OnStop()
	Reset()
	InHost(string) bool
	Filter(string) bool
	Request(req *http.Request) (*http.Request, *http.Response)
}

type ProxyPluginResponse interface {
	OnStart()
	OnStop()
	Reset()
	InHost(string) bool
	Filter(string) bool
	Response(resp *http.Response) *http.Response
}

func PluginLoad(pdir string) {
	files, err := ioutil.ReadDir(pdir)
	if err != nil {
		log.Error(err)
		return
	}
	for _, file := range files {

		fileInfo := strings.Split(file.Name(), ".")
		if !file.IsDir() && len(fileInfo) == 2 && fileInfo[1] == "so" {
			log.Info("Plugin file:", file.Name())
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
			pluginsMap[fileInfo[0]] = plug.(ProxyPlugin)
			pluginsMap[fileInfo[0]].OnStart()
		}
	}

	log.Info(pluginsMap)

}

func PluginOnProxy(proxy *goproxy.ProxyHttpServer) {
	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		for n, _ := range pluginsMap {
			if pp, ok := pluginsMap[n].(ProxyPluginRequest); ok &&
				pp.InHost(req.Host) &&
				pp.Filter(req.URL.Path) {
				log.Info("RUN START Request: ", n)
				res, resp := pp.Request(req)
				return res, resp
			}
		}
		return req, nil
	})

	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		for n, plug := range pluginsMap {
			if pp, ok := plug.(ProxyPluginResponse); ok &&
				resp != nil &&
				plug.InHost(resp.Request.Host) &&
				plug.Filter(resp.Request.RequestURI) {
				log.Info("RUN START Response", n)
				resp = pp.Response(resp)
				return resp
			}
		}
		return resp
	})

}

func PluginOnStop() {
	for n, _ := range pluginsMap {
		//pluginsMap[n].OnStart()
		log.Info("STOP ", n)
		pluginsMap[n].OnStop()
	}
}

func PluginsOnReset() {
	for n, _ := range pluginsMap {
		//pluginsMap[n].OnStart()
		log.Info("RESET ", n)
		pluginsMap[n].Reset()
	}
}
