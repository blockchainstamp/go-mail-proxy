package smtp

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/blockchainstamp/go-mail-proxy/utils/common"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/sirupsen/logrus"
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
		_smtpLog.Error("prepare account failed")
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

	_smtpLog.Info("smtp service init success at:", s.Addr)
	return smtpSrv, nil
}
func (ss *Service) StartWithCtx(cancel context.CancelFunc) error {

	go func() {
		//TODO:: recover the error
		if ss.smtpSrv.AllowInsecureAuth {
			_smtpLog.Info("smtp service running at:", ss.smtpSrv.Addr)
			if err := ss.smtpSrv.ListenAndServe(); err != nil {
				_smtpLog.Warn(err)
			}
		} else {
			_smtpLog.Info("smtp service with tls running at:", ss.smtpSrv.Addr)
			if err := ss.smtpSrv.ListenAndServeTLS(); err != nil {
				_smtpLog.Warn(err)
			}
		}
		cancel()
	}()

	return nil
}
func (ss *Service) Start(sig chan struct{}) error {

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
		sig <- struct{}{}
	}()

	return nil
}

func (ss *Service) Stop() {
	if ss.smtpSrv != nil {
		_ = ss.smtpSrv.Close()
	}
	ss.smtpSrv = nil
}

func (ss *Service) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &Session{
		delegate: ss,
	}, nil
}

func SendMailTls(addr string, auth common.Auth, env *BEnvelope, tls *tls.Config) error {
	a := sasl.NewPlainClient("", auth.UserName, auth.PassWord)
	sender, err := smtp.DialTLS(addr, tls)
	if err != nil {
		_smtpLog.Warnf("dial to %s failed:%s", addr, err)
		return err
	}
	defer sender.Close()
	err = sender.Hello("localhost")
	if err != nil {
		_smtpLog.Warn("hello err:", err)
		return err
	}

	err = sender.Auth(a)
	if err != nil {
		_smtpLog.Warn("auth err:", err)
		return err
	}

	err = sender.Mail(env.From, nil)
	if err != nil {
		_smtpLog.Warn("mail err:", err)
		return err
	}
	for _, to := range env.Tos {
		err = sender.Rcpt(to)
		if err != nil {
			_smtpLog.Warn("rcpt err:", to, err)
			return err
		}
	}
	wc, err := sender.Data()
	if err != nil {
		_smtpLog.Warn("data err:", err)
		return err
	}
	//_, err = io.Copy(wc, env.Data)
	_, err = env.WriteTo(wc)
	if err != nil {
		_smtpLog.Warn("write to err:", err)
		return err
	}
	err = wc.Close()
	if err != nil {
		_smtpLog.Warn("close err:", err)
		return err
	}
	err = sender.Quit()
	if err != nil {
		_smtpLog.Warn("quit err:", err)
		return err
	}
	_smtpLog.Info("send mail success: ", env.From)
	return nil
}
func AuthTls(addr string, auth *common.Auth, tls *tls.Config) error {
	a := sasl.NewPlainClient("", auth.UserName, auth.PassWord)
	sender, err := smtp.DialTLS(addr, tls)
	if err != nil {
		_smtpLog.Warnf("dial to %s failed:%s", addr, err)
		return err
	}
	defer sender.Close()
	err = sender.Hello("localhost")
	if err != nil {
		_smtpLog.Warn("hello err:", err)
		return err
	}

	err = sender.Auth(a)
	if err != nil {
		_smtpLog.Warn("auth err:", err)
		return err
	}
	_smtpLog.Info("tls auth success: ", auth.UserName)
	return nil
}

func AuthNormal(addr string, auth *common.Auth) error {
	a := sasl.NewPlainClient("", auth.UserName, auth.PassWord)
	sender, err := smtp.Dial(addr)
	if err != nil {
		_smtpLog.Warnf("dial to %s failed:%s", addr, err)
		return err
	}
	defer sender.Close()
	err = sender.Hello("localhost")
	if err != nil {
		_smtpLog.Warn("hello err:", err)
		return err
	}
	if ok, _ := sender.Extension("STARTTLS"); !ok {
		_smtpLog.Warn("hello err:", err)
		return fmt.Errorf("smtp: server doesn't support STARTTLS")
	}
	err = sender.StartTLS(nil)
	if err != nil {
		_smtpLog.Warn("hello err:", err)
		return err
	}
	err = sender.Auth(a)
	if err != nil {
		_smtpLog.Warn("auth err:", err)
		return err
	}
	_smtpLog.Info("normal auth success: ", auth.UserName)
	return nil
}

func (ss *Service) SendMail(auth common.Auth, env *BEnvelope) error {
	conf := ss.conf.getRemoteConf(auth.UserName)
	if conf == nil {
		return common.ConfErr
	}
	addr := fmt.Sprintf("%s:%d", conf.RemoteSrvName, conf.RemoteSrvPort)
	return SendMailTls(addr, auth, env, conf.tlsConfig)
}

func (ss *Service) AUTH(auth *common.Auth) error {
	conf := ss.conf.getRemoteConf(auth.UserName)
	if conf == nil {
		return common.ConfErr
	}
	//addr := fmt.Sprintf("%s:%d", conf.RemoteSrvName, conf.RemoteSrvPort)
	//return AuthTls(addr, auth, conf.tlsConfig)
	addr := fmt.Sprintf("%s:%d", conf.RemoteSrvName, 25)
	return AuthNormal(addr, auth)
}
