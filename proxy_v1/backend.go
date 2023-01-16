package proxy_v1

import (
	"crypto/tls"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

var (
	_backendLog = logrus.WithFields(logrus.Fields{
		"mode":    "backend",
		"package": "proxy_v1",
	})
)

type BackendSrv struct {
	tlsCfg *tls.Config
	cfg    *BackendConf
}

func NewBackendServ(conf *BackendConf) (*BackendSrv, error) {
	tlsCfg, err := conf.loadTLSCfg()
	if err != nil {
		return nil, err
	}

	bk := &BackendSrv{
		tlsCfg: tlsCfg,
		cfg:    conf,
	}
	_backendLog.Info("backend service init success")
	return bk, err
}

func (bs *BackendSrv) SendMail(auth Auth, env *BEnvelope) error {
	dialer := gomail.NewDialer(bs.cfg.ServerName, bs.cfg.ServerPort, auth.UserName, auth.PassWord)
	dialer.TLSConfig = bs.tlsCfg

	sender, err := dialer.Dial()
	if err != nil {
		_backendLog.Warnf("dial to %s failed:%s", bs.cfg.ServerName, err)
		return err
	}
	defer sender.Close()
	return sender.Send("ribencong@163.com", []string{"ribencong@126.com"}, env)
}
