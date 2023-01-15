package proxy_v1

import (
	"os"
)

type ProxyService struct {
	backend *BackendSrv
	smtp    *SMTPSrv
}

func (p *ProxyService) InitByConf(confPath string) error {
	if err := _srvConf.active(confPath); err != nil {
		return err
	}

	bk, err := NewBackendServ(_srvConf.BackendConf)
	if err != nil {
		return err
	}

	smtp, err := NewSMTPSrv(_srvConf.SMTPConf)
	if err != nil {
		return err
	}

	p.smtp = smtp
	p.backend = bk
	return nil
}

func (p *ProxyService) Start() error {
	return nil
}

func (p *ProxyService) ShutDown() error {
	os.Exit(0)
	return nil
}

func NewProxy() *ProxyService {
	ps := &ProxyService{}
	return ps
}