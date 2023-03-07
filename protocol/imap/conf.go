package imap

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/blockchainstamp/go-mail-proxy/utils/common"
	"os"
	"strings"
)

type RemoteConf struct {
	RemoteSrvCAs   string `json:"ca_files"`
	RemoteCADomain string `json:"ca_domain"`
	AllowNotSecure bool   `json:"allow_not_secure"`
	RemoteSrvName  string `json:"remote_srv_name"`
	RemoteSrvPort  int    `json:"remote_srv_port"`
	tlsConfig      *tls.Config
	remoteSrvAddr  string
}

func (rc *RemoteConf) String() string {
	s := "\nRoot CAs:\t" + rc.RemoteSrvCAs
	s += "\nCA Domain:\t" + rc.RemoteCADomain
	s += fmt.Sprintf("\nAllow not security:\t%t", rc.AllowNotSecure)
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
		conf.remoteSrvAddr = fmt.Sprintf("%s:%d", conf.RemoteSrvName, conf.RemoteSrvPort)

		if conf.AllowNotSecure {
			_imapLog.Info("no need ca file:", conf.RemoteSrvCAs)
			continue
		}
		fileNames := strings.Split(conf.RemoteSrvCAs, common.CAFileSep)
		if len(fileNames) == 0 {
			_imapLog.Errorf("no valid ca file:[%s]", conf.RemoteSrvCAs)
			return common.TLSErr
		}
		rootCAs := x509.NewCertPool()
		for _, caPath := range fileNames {
			_imapLog.Debug("ca file path:", caPath)
			data, err := os.ReadFile(caPath)
			if err != nil {
				_imapLog.Errorf("read ca file[%s] failed:%s", caPath, err)
				return err
			}
			rootCAs.AppendCertsFromPEM(data)
		}
		tlsConf := &tls.Config{
			ServerName: conf.RemoteCADomain,
			RootCAs:    rootCAs,
		}
		conf.tlsConfig = tlsConf
	}
	return nil
}

func (c *Conf) getRemoteConf(mailAddr string) *RemoteConf {
	cfg, ok := c.RemoteConf[mailAddr]
	if ok {
		return cfg
	}
	var addr = strings.Split(mailAddr, common.MailAddrSep)
	if len(addr) != 2 {
		_imapLog.Warn("invalid email address:", mailAddr)
		return nil
	}
	var domain = addr[1]
	return c.RemoteConf[domain]
}
