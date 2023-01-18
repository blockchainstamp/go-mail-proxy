package proxy_v1

import (
	"crypto/tls"
	"fmt"
	"github.com/blockchainstamp/go-mail-proxy/utils"
	"github.com/sirupsen/logrus"
)

var (
	_srvConf = &Config{}
)

type Config struct {
	LogLevel          uint32 `json:"log_level"`
	SMTPConfPath      string `json:"smtp"`
	IMAPConfPath      string `json:"imap"`
	AllowInsecureAuth bool   `json:"allow-insecure-auth"`
	TlsKeyPath        string `json:"tls-key-path"`
	TlsCertPath       string `json:"tls-cert-path"`
}

func (c *Config) String() string {
	s := "\n+++++++++++++++++++++++config+++++++++++++++++++++++++++++"
	s += "\nLog Level:\t" + logrus.Level(c.LogLevel).String()
	s += "\nSMTP Config:\t" + c.SMTPConfPath
	s += "\nIMAP Config:\t" + c.IMAPConfPath
	s += fmt.Sprintf("\nSecure Auth:\t%t", c.AllowInsecureAuth)
	s += "\nTls Key Path:\t" + c.TlsKeyPath
	s += "\nTls Cert Path:\t" + c.TlsCertPath
	s += "\n++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
	return s
}

func (c *Config) prepare(confPath string) error {
	var (
		err error = nil
	)

	if err = utils.ReadJsonFile(confPath, _srvConf); err != nil {
		return err
	}
	fmt.Println(_srvConf.String())

	logrus.SetLevel(logrus.Level(_srvConf.LogLevel))
	logrus.SetFormatter(&logrus.JSONFormatter{})

	return err
}

func (c *Config) loadServerTlsCnf() (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(c.TlsCertPath, c.TlsKeyPath)
	if err != nil {
		return nil, err
	}
	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}

	return cfg, err
}
