// Able Gao @
// ablegao@gmail.com
// description：
//
//

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"yfiddler/certs"
	"yfiddler/header"

	"github.com/elazarl/goproxy"
	log "github.com/sirupsen/logrus"
)

var (
	proxy     *goproxy.ProxyHttpServer
	verbose   = flag.Bool("v", false, "should every proxy request be logged to stdout")
	addr      = flag.String("addr", ":8080", "proxy listen address")
	confFile  = flag.String("config", "./configure/config.yaml", "Configure file path . ")
	pluginDir = flag.String("plugins", "./plugins", "Plugins Folder path.")
)

func writeToFile(p string, value string) {
	f, err := os.OpenFile(p, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	f.WriteString(value)
}

func main() {
	writeToFile("./reload.sh", fmt.Sprintf("kill -SIGUSR1 %d", os.Getpid()))
	flag.Parse()
	certs.UpdateCA()
	header.PluginLoad(*pluginDir)
	proxy = goproxy.NewProxyHttpServer()
	proxy.Verbose = *verbose
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	header.YamlHeader(proxy, *confFile)
	header.PluginOnProxy(proxy)
	log.Info("load compiled!")

	go func() {
		err := http.ListenAndServe(*addr, proxy)
		log.Error(err)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt, syscall.SIGUSR1)

	for {
		sig := <-sigs

		switch sig {
		case syscall.SIGKILL, os.Interrupt, syscall.SIGTERM, syscall.SIGINT:
			return
		case syscall.SIGUSR1:
			//ReloadPackageList()
		default:
			fmt.Println(sig)
		}

	}
}