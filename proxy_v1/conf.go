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
	_srvConf = &Config{
		SMTPConf:    &SMTPConf{},
		BackendConf: &BackendConf{},
	}
	TLSErr = errors.New("no valid tls config")
)

type Config struct {
	LogLevel        uint32 `json:"log_level"`
	SMTPConfPath    string `json:"smtp"`
	BackendConfPath string `json:"backend"`
	*SMTPConf       `json:"-"`
	*BackendConf    `json:"-"`
}

func (c *Config) String() string {
	s := "\n+++++++++++++++++++++++config+++++++++++++++++++++++++++++"
	s += "\nLog Level:\t" + logrus.Level(c.LogLevel).String()
	s += "\nSMTP Config:\t" + c.SMTPConfPath
	s += "\nBackend Config:\t" + c.BackendConfPath
	s += "\n++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
	return s
}

func (c *Config) prepare(confPath string) error {
	var (
		err error = nil
	)

	if err = prepareConf(confPath, _srvConf); err != nil {
		return err
	}
	fmt.Println(_srvConf.String())

	if err = prepareConf(_srvConf.SMTPConfPath, _srvConf.SMTPConf); err != nil {
		return err
	}
	fmt.Println(_srvConf.SMTPConf.String())

	if err = prepareConf(_srvConf.BackendConfPath, _srvConf.BackendConf); err != nil {
		return err
	}
	fmt.Println(_srvConf.BackendConf.String())

	logrus.SetLevel(logrus.Level(_srvConf.LogLevel))
	logrus.SetFormatter(&logrus.JSONFormatter{})

	return err
}

func prepareConf(confPath string, conf interface{}) error {
	data, err := os.ReadFile(confPath)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, conf); err != nil {
		return err
	}
	return nil
}

type SMTPConf struct {
}

func (sc *SMTPConf) String() string {
	s := "\n=========service============="
	s += "\n=============================\n"
	return s
}

type BackendConf struct {
	RootCAFiles string `json:"ca_files"`
	ServerName  string `json:"server_name"`
	ServerPort  int    `json:"server_port"`
}

func (bc *BackendConf) String() string {
	s := "\n==========backend============"
	s += "\nRoot CAs:\t" + bc.RootCAFiles
	s += "\nServer Name:\t" + bc.ServerName
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
