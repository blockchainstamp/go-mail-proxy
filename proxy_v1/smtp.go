package proxy_v1

import "github.com/sirupsen/logrus"

var (
	_smptLog = logrus.WithFields(logrus.Fields{
		"mode":    "smtp service",
		"package": "proxy_v1",
	})
)

type SMTPSrv struct {
}

func NewSMTPSrv(conf *SMTPConf) (*SMTPSrv, error) {
	ss := &SMTPSrv{}
	_backendLog.Info("smtp receiving service init success")
	return ss, nil
}
