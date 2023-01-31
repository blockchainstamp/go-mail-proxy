package smtp

import (
	"crypto/tls"
	common2 "github.com/blockchainstamp/go-mail-proxy/protocol/common"
	"github.com/emersion/go-smtp"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"time"
)

var (
	_smtpLog = logrus.WithFields(logrus.Fields{
		"mode":    "smtp service",
		"package": "smtp",
	})
)

type Service struct {
	smtpSrv *smtp.Server
	conf    *Conf
}

func NewSMTPSrv(conf *Conf, lclSrvTls *tls.Config) (*Service, error) {
	if err := conf.prepareAccounts(); err != nil {
		return nil, err
	}

	smtpSrv := &Service{
		conf: conf,
	}
	s := smtp.NewServer(smtpSrv)

	s.Addr = conf.SrvAddr
	s.Domain = conf.SrvDomain
	s.ReadTimeout = time.Duration(conf.ReadTimeOut) * time.Second
	s.WriteTimeout = time.Duration(conf.WriteTimeOut) * time.Second
	s.MaxMessageBytes = conf.MaxMessageBytes
	s.MaxRecipients = conf.MaxRecipients
	s.AllowInsecureAuth = lclSrvTls == nil
	s.TLSConfig = lclSrvTls

	smtpSrv.smtpSrv = s

	_smtpLog.Info("smtp receiving service init success at:", s.Addr)
	return smtpSrv, nil
}

func (ss *Service) Start() error {

	go func() {
		//TODO:: recover the error
		if ss.smtpSrv.AllowInsecureAuth {
			_smtpLog.Info("smtp service running at:", ss.smtpSrv.Addr)
			if err := ss.smtpSrv.ListenAndServe(); err != nil {
				panic(err)
			}
		} else {
			_smtpLog.Info("smtp service with tls running at:", ss.smtpSrv.Addr)
			if err := ss.smtpSrv.ListenAndServeTLS(); err != nil {
				panic(err)
			}
		}
	}()

	return nil
}

func (ss *Service) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{delegate: ss, env: &BEnvelope{}}, nil
}

func (ss *Service) SendMail(auth common2.Auth, env *BEnvelope) error {
	conf := ss.conf.getRemoteConf(auth.UserName)
	if conf == nil {
		return common2.ConfErr
	}
	dialer := gomail.NewDialer(conf.RemoteSrvName, conf.RemoteSrvPort, auth.UserName, auth.PassWord)
	dialer.TLSConfig = conf.tlsConfig

	sender, err := dialer.Dial()
	if err != nil {
		_smtpLog.Warnf("dial to %s failed:%s", conf.RemoteSrvName, err)
		return err
	}
	defer sender.Close()
	err = sender.Send(env.From, env.Tos, env)
	if err != nil {
		_smtpLog.Warnf("SendMail failed :%s", err)
	}
	return err
}

func (ss *Service) AUTH(auth *common2.Auth) error {
	conf := ss.conf.getRemoteConf(auth.UserName)
	if conf == nil {
		return common2.ConfErr
	}
	dialer := gomail.NewDialer(conf.RemoteSrvName, conf.RemoteSrvPort, auth.UserName, auth.PassWord)
	dialer.TLSConfig = conf.tlsConfig

	sender, err := dialer.Dial()
	if err != nil {
		_smtpLog.Warnf("dial to %s failed:%s", conf.RemoteSrvName, err)
		return err
	}
	_smtpLog.Info("auth success:", auth.UserName)
	return sender.Close()
}
