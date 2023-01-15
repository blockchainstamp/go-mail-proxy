package proxy_v1

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

const (
	CAFileSep = ";"
)

var (
	_srvConf = &Config{}
	TLSErr   = errors.New("no valid tls config")
)

type Config struct {
	LogLevel        uint32 `json:"log_level"`
	SMTPConfPath    string `json:"smtp"`
	BackendConfPath string `json:"backend"`
	*SMTPConf       `json:"-"`
	*BackendConf    `json:"-"`
}

func (c *Config) String() string {
	s := "\n++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
	s += "Log Level:\t" + logrus.Level(c.LogLevel).String()
	s += "SMTP Config:\t" + c.SMTPConfPath
	s += "Backend Config:\t" + c.BackendConfPath
	s += "\n++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
	return s
}

func (c *Config) active(confPath string) error {
	var (
		data []byte = nil
		err  error  = nil
	)

	data, err = os.ReadFile(confPath)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, _srvConf); err != nil {
		return err
	}

	fmt.Println(_srvConf.String())

	logrus.SetLevel(logrus.Level(_srvConf.LogLevel))

	if err = _srvConf.activeSMTPConf(); err != nil {
		return err
	}
	if err = _srvConf.activeBackendConf(); err != nil {
		return err
	}
	return err
}

func (c *Config) activeSMTPConf() error {
	data, err := os.ReadFile(_srvConf.SMTPConfPath)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, _srvConf.SMTPConf); err != nil {
		return err
	}

	fmt.Println(_srvConf.SMTPConf.String())
	return nil
}

func (c *Config) activeBackendConf() error {
	data, err := os.ReadFile(_srvConf.BackendConfPath)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, _srvConf.BackendConf); err != nil {
		return err
	}
	fmt.Println(_srvConf.BackendConf.String())
	return nil
}

type SMTPConf struct {
}

func (sc *SMTPConf) String() string {
	s := "\n=========service=============\n"
	s += "\n=============================\n"
	return s
}

type BackendConf struct {
	RootCAFiles string `json:"ca_files"`
	ServerName  string `json:"server_name"`
}

func (bc *BackendConf) String() string {
	s := "\n==========backend============\n"
	s += "Root CAs:\t" + bc.RootCAFiles
	s += "Server Name:\t" + bc.ServerName
	s += "\n=============================\n"
	return s
}

func (bc *BackendConf) loadTLSCfg() (*tls.Config, error) {

	fileNames := strings.Split(bc.RootCAFiles, CAFileSep)
	if len(fileNames) == 0 {
		return nil, TLSErr
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
		ServerName: bc.ServerName,
		RootCAs:    rootCAs,
	}
	return tlsConf, nil
}
