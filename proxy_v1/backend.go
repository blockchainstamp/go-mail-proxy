package proxy_v1

import (
	"crypto/tls"
	"github.com/emersion/go-imap/backend/memory"
	"github.com/emersion/go-imap/server"
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
	remoteSmtpTls *tls.Config
	cfg           *BackendConf
	imap          *server.Server
}

func newImap() {

}
func NewBackendServ(conf *BackendConf, srvTlsCfg *tls.Config) (*BackendSrv, error) {
	remoteSmtTls, err := conf.loadRemoteRootCAs()
	if err != nil {
		return nil, err
	}

	be := memory.New()

	s := server.New(be)
	s.Addr = conf.ImapAddr
	s.AllowInsecureAuth = srvTlsCfg == nil
	s.TLSConfig = srvTlsCfg

	bk := &BackendSrv{
		remoteSmtpTls: remoteSmtTls,
		cfg:           conf,
		imap:          s,
	}
	_backendLog.Info("backend service init success imap:", s.Addr)
	return bk, err
}

func (bs *BackendSrv) SendMail(auth Auth, env *BEnvelope) error {
	dialer := gomail.NewDialer(bs.cfg.ServerName, bs.cfg.ServerPort, auth.UserName, auth.PassWord)
	dialer.TLSConfig = bs.remoteSmtpTls

	sender, err := dialer.Dial()
	if err != nil {
		_backendLog.Warnf("dial to %s failed:%s", bs.cfg.ServerName, err)
		return err
	}
	defer sender.Close()
	return sender.Send(env.From, env.Tos, env)
}

func (bs *BackendSrv) Start() error {
	go func() {
		if bs.imap.AllowInsecureAuth {
			_backendLog.Info("backend imap start success")
			if err := bs.imap.ListenAndServe(); err != nil {
				panic(err)
			}
		} else {
			_backendLog.Info("backend imap with tls start success")
			if err := bs.imap.ListenAndServeTLS(); err != nil {
				panic(err)
			}
		}
	}()
	return nil
}
