package imap

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/blockchainstamp/go-mail-proxy/proxy_v1/common"
	"os"
	"strings"
)

type RemoteConf struct {
	RemoteSrvCAs  string `json:"ca_files"`
	RemoteSrvName string `json:"remote_srv_name"`
	RemoteSrvPort int    `json:"remote_srv_port"`
	tlsConfig     *tls.Config
	remoteSrvAddr string
}

func (rc *RemoteConf) String() string {
	s := "\nRoot CAs:\t" + rc.RemoteSrvCAs
	s += "\nRemote Server:\t" + rc.RemoteSrvName
	s += fmt.Sprintf("\nRemote Port:\t%d", rc.RemoteSrvPort)
	return s
}

type Conf struct {
	SrvAddr    string                 `json:"srv_addr"`
	SrvDomain  string                 `json:"srv_domain"`
	RemoteConf map[string]*RemoteConf `json:"remote_conf"`
}

func (c *Conf) String() string {
	s := "\n=========service[imap]============="
	s += "\nServer Addr:\t" + c.SrvAddr
	s += "\nServer Domain:\t" + c.SrvDomain
	for r, conf := range c.RemoteConf {
		s += fmt.Sprintf("\n---%s---", r)
		s += conf.String()
		s += fmt.Sprintf("\n------------")
	}
	s += "\n=============================\n"
	return s
}

func (c *Conf) loadRemoteRootCAs() error {
	for _, conf := range c.RemoteConf {
		fileNames := strings.Split(conf.RemoteSrvCAs, common.CAFileSep)
		if len(fileNames) == 0 {
			return common.TLSErr
		}
		rootCAs := x509.NewCertPool()
		for _, caPath := range fileNames {
			data, err := os.ReadFile(caPath)
			if err != nil {
				return err
			}
			rootCAs.AppendCertsFromPEM(data)
		}
		tlsConf := &tls.Config{
			ServerName: conf.RemoteSrvName,
			RootCAs:    rootCAs,
		}
		conf.tlsConfig = tlsConf
		conf.remoteSrvAddr = fmt.Sprintf("%s:%d", conf.RemoteSrvName, conf.RemoteSrvPort)
	}
	return nil
}

func (c *Conf) getRemoteConf(user string) *RemoteConf {
	return c.RemoteConf[user]
}
