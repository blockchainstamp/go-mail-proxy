package proxy_v1

import (
	"github.com/emersion/go-smtp"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	_smptLog = logrus.WithFields(logrus.Fields{
		"mode":    "smtp service",
		"package": "proxy_v1",
	})
)

type SMTPSrv struct {
	smtpSrv *smtp.Server
	conf    *SMTPConf
}

func NewSMTPSrv(conf *SMTPConf, be smtp.Backend) (*SMTPSrv, error) {

	s := smtp.NewServer(be)

	s.Addr = conf.Addr
	s.Domain = conf.Domain
	s.ReadTimeout = time.Duration(conf.ReadTimeOut) * time.Second
	s.WriteTimeout = time.Duration(conf.WriteTimeOut) * time.Second
	s.MaxMessageBytes = conf.MaxMessageBytes
	s.MaxRecipients = conf.MaxRecipients
	s.AllowInsecureAuth = conf.AllowInsecureAuth

	if !conf.AllowInsecureAuth {
		cfg, err := conf.loadServerTlsCnf()
		if err != nil {
			return nil, err
		}
		s.TLSConfig = cfg
	}

	smtpSrv := &SMTPSrv{
		smtpSrv: s,
		conf:    conf,
	}
	_smptLog.Info("smtp receiving service init success at:", s.Addr)
	return smtpSrv, nil
}

func (ss *SMTPSrv) Start() error {

	go func() {
		//TODO:: recover the error
		if ss.smtpSrv.AllowInsecureAuth {
			_smptLog.Info("smtp service running at:", ss.smtpSrv.Addr)
			if err := ss.smtpSrv.ListenAndServe(); err != nil {
				panic(err)
			}
		} else {
			_smptLog.Info("smtp service with tls running at:", ss.smtpSrv.Addr)
			if err := ss.smtpSrv.ListenAndServeTLS(); err != nil {
				panic(err)
			}
		}
	}()

	return nil
}
