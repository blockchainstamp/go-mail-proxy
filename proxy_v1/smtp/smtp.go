package smtp

import (
	"crypto/tls"
	"github.com/blockchainstamp/go-mail-proxy/proxy_v1/common"
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
	smtpSrv      *smtp.Server
	conf         *Conf
	remoteTlsCfg *tls.Config
}

func NewSMTPSrv(conf *Conf, lclSrvTls *tls.Config) (*Service, error) {
	remoteSmtTls, err := conf.loadRemoteRootCAs()
	if err != nil {
		return nil, err
	}

	smtpSrv := &Service{
		conf:         conf,
		remoteTlsCfg: remoteSmtTls,
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
	return &Session{delegate: ss, env: &common.BEnvelope{}}, nil
}

func (ss *Service) SendMail(auth common.Auth, env *common.BEnvelope) error {
	dialer := gomail.NewDialer(ss.conf.RemoteSrvName, ss.conf.RemoteSrvPort, auth.UserName, auth.PassWord)
	dialer.TLSConfig = ss.remoteTlsCfg

	sender, err := dialer.Dial()
	if err != nil {
		_smtpLog.Warnf("dial to %s failed:%s", ss.conf.RemoteSrvName, err)
		return err
	}
	defer sender.Close()
	return sender.Send(env.From, env.Tos, env)
}

func (ss *Service) AUTH(auth *common.Auth) error {
	dialer := gomail.NewDialer(ss.conf.RemoteSrvName, ss.conf.RemoteSrvPort, auth.UserName, auth.PassWord)
	dialer.TLSConfig = ss.remoteTlsCfg

	sender, err := dialer.Dial()
	if err != nil {
		_smtpLog.Warnf("dial to %s failed:%s", ss.conf.RemoteSrvName, err)
		return err
	}
	sender.Close()
	return nil
}
