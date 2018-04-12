package certs

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"

	"github.com/elazarl/goproxy"
)

var (
	caCert, caKey []byte
	certPem       = flag.String("cert", "./rootCA.pem", "")
	certKey       = flag.String("certkey", "./rootCA.key", "")
)

func UpdateCA() error {
	var err error
	caCert, err = ioutil.ReadFile(*certPem)
	if err != nil {
		return err
	}
	caKey, err = ioutil.ReadFile(*certKey)
	if err != nil {
		return err
	}
	goproxyCa, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		return err
	}
	if goproxyCa.Leaf, err = x509.ParseCertificate(goproxyCa.Certificate[0]); err != nil {
		return err
	}
	goproxy.GoproxyCa = goproxyCa
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	return nil
}
