package header

import (
	"net/http"
	"strings"
	"yfiddler/configure"
	"yfiddler/utils"

	"github.com/elazarl/goproxy"
	log "github.com/sirupsen/logrus"
)

func YamlHeader(proxy *goproxy.ProxyHttpServer, yamlFile string) {
	if yamlFile == "" {
		return
	}
	config, err := configure.NewConfigureByYaml(yamlFile)
	if err != nil || config == nil {
		log.Error(err)
		return
	}
	for _, r := range config.Requests {
		log.Debug("YAML: ", utils.Json(r))
	}

	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		for id, request := range config.Requests {

			if config.Requests[id].InHosts(req.Host) && request.Filter(req.RequestURI) {
				log.Debug("RUN YAML Request:", request, " true")
				for hk, hv := range request.Headers {
					req.Header.Set(hk, strings.Join(hv, "; "))
				}
				if request.Response != nil {
					res := &http.Response{}
					res.Request = req
					res.TransferEncoding = req.TransferEncoding
					res.Header = http.Header(request.Headers)
					size, b, err := request.Response.GenReadCloser()
					if err != nil {
						log.Error(err)
						return req, nil
					}
					res.Body = b
					res.ContentLength = size
					return req, res
				}
				log.Error("Reqeuset filter error:", req.URL.String())
				return req, nil
			} else {
				log.Debug("YAML Filter OUT: ", request.InHosts(req.Host), " ", request.Filter(req.RequestURI), request)
			}

		}
		return req, nil
	})

	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		var err error
		for n, res := range config.Responses {
			if resp != nil && res.InHosts(resp.Request.Host) && res.Filter(resp.Request.RequestURI) {
				if res.DataType == configure.DATA_TYPE_NONE {
					for hk, hv := range res.Headers {
						resp.Header.Set(hk, strings.Join(hv, "; "))
					}
				} else {

					resp.ContentLength, resp.Body, err = config.Responses[n].GenReadCloser()
					if err != nil {
						log.Error(err)
						return resp
					}
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
