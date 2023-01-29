package proxy_v1

import (
	"crypto/tls"
	"fmt"
	"github.com/blockchainstamp/go-mail-proxy/proxy_v1/common"
	"github.com/blockchainstamp/go-mail-proxy/proxy_v1/imap"
	"github.com/blockchainstamp/go-mail-proxy/proxy_v1/smtp"
	"github.com/blockchainstamp/go-mail-proxy/utils"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	_proxyLog = logrus.WithFields(logrus.Fields{
		"mode":    "proxy main process",
		"package": "proxy",
	})
)

type ProxyService struct {
	imapSrv *imap.Service
	smtpSrv *smtp.Service
}

func (p *ProxyService) InitByConf(confPath string) error {
	if err := _srvConf.prepare(confPath); err != nil {
		return err
	}

	var localTlsCfg *tls.Config
	if !_srvConf.AllowInsecureAuth {
		cfg, err := _srvConf.loadServerTlsCnf()
		if err != nil {
			return err
		}
		localTlsCfg = cfg
	}
	smtpCfg := &smtp.Conf{}
	if err := utils.ReadJsonFile(_srvConf.SMTPConfPath, smtpCfg); err != nil {
		return err
	}
	fmt.Println(smtpCfg.String())

	imapCfg := &imap.Conf{}
	if err := utils.ReadJsonFile(_srvConf.IMAPConfPath, imapCfg); err != nil {
		return err
	}
	fmt.Println(imapCfg.String())

	is, err := imap.NewIMAPSrv(imapCfg, localTlsCfg)
	if err != nil {
		return err
	}

	ss, err := smtp.NewSMTPSrv(smtpCfg, localTlsCfg)
	if err != nil {
		return err
	}
	go utils.StartCmdService(_srvConf.CmdSrvAddr)
	p.smtpSrv = ss
	p.imapSrv = is
	_proxyLog.Info("proxy process init success")
	return nil
}

func (p *ProxyService) Start() error {
	var err error = nil
	if err = p.imapSrv.Start(); err != nil {
		return err
	}
	if err = p.smtpSrv.Start(); err != nil {
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

func (p *ProxyService) Command(cmd common.Command) any {
	return nil
}

func NewProxy() *ProxyService {
	ps := &ProxyService{}
	common.RegCmdProc(common.CMDProxy, ps.Command)
	return ps
}
