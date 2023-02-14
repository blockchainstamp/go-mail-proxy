package proxy_v1

import (
	"crypto/tls"
	"fmt"
	"github.com/blockchainstamp/go-mail-proxy/protocol/common"
	imap2 "github.com/blockchainstamp/go-mail-proxy/protocol/imap"
	smtp2 "github.com/blockchainstamp/go-mail-proxy/protocol/smtp"
	"github.com/blockchainstamp/go-mail-proxy/utils"
	bstamp "github.com/blockchainstamp/go-stamp-wallet"
	"github.com/blockchainstamp/go-stamp-wallet/comm"
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
	imapSrv *imap2.Service
	smtpSrv *smtp2.Service
}

func (p *ProxyService) InitByConf(conf any, auth string) error {
	_srvConf = conf.(*Config)
	fmt.Println(_srvConf.String())
	if err := bstamp.InitSDK(_srvConf.StampDBPath); err != nil {
		return err
	}
	level, err := logrus.ParseLevel(_srvConf.LogLevel)
	if err != nil {
		return err
	}
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	var localTlsCfg *tls.Config
	if !_srvConf.AllowInsecureAuth {
		cfg, err := _srvConf.loadServerTlsCnf()
		if err != nil {
			return err
		}
		localTlsCfg = cfg
	}
	smtpCfg := &smtp2.Conf{}
	if err := utils.ReadJsonFile(_srvConf.SMTPConfPath, smtpCfg); err != nil {
		return err
	}
	fmt.Println(smtpCfg.String())

	imapCfg := &imap2.Conf{}
	if err := utils.ReadJsonFile(_srvConf.IMAPConfPath, imapCfg); err != nil {
		return err
	}
	fmt.Println(imapCfg.String())

	is, err := imap2.NewIMAPSrv(imapCfg, localTlsCfg)
	if err != nil {
		return err
	}

	if len(auth) > 0 && len(smtpCfg.StampWalletAddr) > 0 {
		w, err := bstamp.Inst().ActiveWallet(comm.WalletAddr(smtpCfg.StampWalletAddr), auth)
		if err != nil {
			return err
		}
		fmt.Println("wallet address:", w.Address())
		fmt.Println("eth address:", w.EthAddr())
	}

	ss, err := smtp2.NewSMTPSrv(smtpCfg, localTlsCfg)
	if err != nil {
		return err
	}
	go utils.StartCmdService(_srvConf.CmdSrvAddr, nil)
	p.smtpSrv = ss
	p.imapSrv = is
	_proxyLog.Info("proxy process init success")
	return nil
}

func (p *ProxyService) StartWithSig(sig chan struct{}) error {
	var err error = nil
	if err = p.imapSrv.Start(sig); err != nil {
		return err
	}
	if err = p.smtpSrv.Start(sig); err != nil {
		return err
	}
	_proxyLog.Info("proxy process start success")
	go p.monitor(sig)
	return nil
}

func (p *ProxyService) Start() error {
	var sig = make(chan struct{}, 2)
	return p.StartWithSig(sig)
}

func (p *ProxyService) monitor(sig chan struct{}) {
	for {
		select {
		case <-sig:
			_ = p.ShutDown()
		}
	}
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
