package proxy_v1

import (
	"github.com/emersion/go-smtp"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	_proxyLog = logrus.WithFields(logrus.Fields{
		"mode":    "proxy main process",
		"package": "proxy_v1",
	})
)

type ProxyService struct {
	backend *BackendSrv
	smtp    *SMTPSrv
}

func (p *ProxyService) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{sender: p.backend, env: &BEnvelope{}}, nil
}

func (p *ProxyService) InitByConf(confPath string) error {
	if err := _srvConf.prepare(confPath); err != nil {
		return err
	}

	bk, err := NewBackendServ(_srvConf.BackendConf)
	if err != nil {
		return err
	}

	ss, err := NewSMTPSrv(_srvConf.SMTPConf, p)
	if err != nil {
		return err
	}

	p.smtp = ss
	p.backend = bk
	_proxyLog.Info("proxy process init success")
	return nil
}

func (p *ProxyService) Start() error {
	var err error = nil
	if err = p.backend.Start(); err != nil {
		return err
	}
	if err = p.smtp.Start(); err != nil {
		return err
	}
	_proxyLog.Info("proxy process start success")

	return nil
}

func (p *ProxyService) ShutDown() error {
	_proxyLog.Info("proxy process shutdown")
	// TODO::
	os.Exit(0)
	return nil
}

func NewProxy() *ProxyService {
	ps := &ProxyService{}
	return ps
}
