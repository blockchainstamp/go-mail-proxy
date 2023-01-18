package imap

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/blockchainstamp/go-mail-proxy/proxy_v1/common"
	"os"
	"strings"
)

type Conf struct {
	SrvAddr   string `json:"srv_addr"`
	SrvDomain string `json:"srv_domain"`

	RemoteSrvCAs  string `json:"ca_files"`
	RemoteSrvName string `json:"remote_srv_name"`
	RemoteSrvPort int    `json:"remote_srv_port"`
}

func (c *Conf) String() string {
	s := "\n=========service============="
	s += "\nServer Addr:\t" + c.SrvAddr
	s += "\nServer Domain:\t" + c.SrvDomain
	s += "\nRoot CAs:\t" + c.RemoteSrvCAs
	s += "\nRemote Server:\t" + c.RemoteSrvName
	s += fmt.Sprintf("\nRemote Port:\t%d", c.RemoteSrvPort)
	s += "\n=============================\n"
	return s
	return s
}
func (c *Conf) loadRemoteRootCAs() (*tls.Config, error) {

	fileNames := strings.Split(c.RemoteSrvCAs, common.CAFileSep)
	if len(fileNames) == 0 {
		return nil, common.TLSErr
	}
	rootCAs := x509.NewCertPool()
	for _, caPath := range fileNames {
		data, err := os.ReadFile(caPath)
		if err != nil {
			return nil, err
		}
		rootCAs.AppendCertsFromPEM(data)
	}
	tlsConf := &tls.Config{
		ServerName: c.RemoteSrvName,
		RootCAs:    rootCAs,
	}
	return tlsConf, nil
}
