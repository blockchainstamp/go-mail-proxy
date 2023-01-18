package imap

import (
	"crypto/tls"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
	"github.com/emersion/go-imap/server"
	"github.com/sirupsen/logrus"
)

var (
	_imapLog = logrus.WithFields(logrus.Fields{
		"mode":    "smtp service",
		"package": "imap",
	})
)

type Service struct {
	users        map[string]*User
	remoteTlsCfg *tls.Config
	srv          *server.Server
}

func (is *Service) Login(_ *imap.ConnInfo, username, password string) (backend.User, error) {

	u := &User{username: username, password: password}

	u.mailboxes = map[string]*Mailbox{
		"INBOX": {
			name: "INBOX",
			user: u,
		},
	}
	return u, nil
}

func NewIMAPSrv(cfg *Conf, lclSrvTls *tls.Config) (*Service, error) {

	remoteSmtTls, err := cfg.loadRemoteRootCAs()
	if err != nil {
		return nil, err
	}
	is := &Service{
		remoteTlsCfg: remoteSmtTls,
		users:        make(map[string]*User),
	}
	s := server.New(is)
	s.Addr = cfg.SrvAddr
	s.AllowInsecureAuth = lclSrvTls == nil
	s.TLSConfig = lclSrvTls

	is.srv = s
	return is, nil
}

func (is *Service) Start() error {
	go func() {
		if is.srv.AllowInsecureAuth {
			_imapLog.Info("backend imap start success")
			if err := is.srv.ListenAndServe(); err != nil {
				panic(err)
			}
		} else {
			_imapLog.Info("backend imap with tls start success")
			if err := is.srv.ListenAndServeTLS(); err != nil {
				panic(err)
			}
		}
	}()
	return nil
}
