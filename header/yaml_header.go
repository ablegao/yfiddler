package header

import (
	"net/http"
	"yfiddler/configure"

	"github.com/elazarl/goproxy"
	log "github.com/sirupsen/logrus"
)

func YamlHeader(proxy *goproxy.ProxyHttpServer, yamlFile string) {
	if yamlFile == "" {
		return
	}
	config, err := configure.NewConfigureByYaml(yamlFile)
	if err != nil {
		log.Error(err)
		return
	}

	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {

		for _, request := range config.Requests {

			if request.InHosts(req.Host) && request.Filter(req.RequestURI) {
				for hk, hv := range request.Headers {
					req.Header.Set(hk, hv)
				}
				if request.Response != nil {
					res := &http.Response{}
					if request.Response.GetResponse(req, res) {
						return req, res
					}
				}
				return req, nil
			}

		}
		return req, nil
	})

	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		for _, res := range config.Responses {
			if resp.Request != nil && res.InHosts(resp.Request.Host) && res.Filter(resp.Request.RequestURI) {
				if res.DataType == configure.DATA_TYPE_NONE {
					for hk, hv := range res.Headers {
						resp.Header.Set(hk, hv)
					}
				} else {
					res.GetResponse(resp.Request, resp)
				}
				return resp
			}
		}
		return resp
	})

}
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
