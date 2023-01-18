package imap

import (
	"crypto/tls"
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-imap/server"
	"github.com/sirupsen/logrus"
)

var (
	_imapLog = logrus.WithFields(logrus.Fields{
		"mode":    "imap service",
		"package": "imap",
	})
)

type Service struct {
	users         map[string]*User
	remoteTlsCfg  *tls.Config
	remoteSrvAddr string
	srv           *server.Server
}

func (is *Service) Login(_ *imap.ConnInfo, username, password string) (backend.User, error) {

	u := &User{username: username, password: password}
	c, err := client.DialTLS(is.remoteSrvAddr, is.remoteTlsCfg)
	if err != nil {
		_imapLog.Warn("dial failed", is.remoteSrvAddr, err)
		return nil, err
	}

	defer c.Logout()

	if err := c.Login(username, password); err != nil {
		_imapLog.Warnf("user[%s] login failed:%s", username, err)
		return nil, err
	}

	_ = u.CreateMailbox("INBOX")
	is.users[username] = u
	_imapLog.Infof("user[%s] login success", username)
	return u, nil
}

func NewIMAPSrv(cfg *Conf, lclSrvTls *tls.Config) (*Service, error) {

	remoteSmtTls, err := cfg.loadRemoteRootCAs()
	if err != nil {
		return nil, err
	}
	is := &Service{
		remoteTlsCfg:  remoteSmtTls,
		users:         make(map[string]*User),
		remoteSrvAddr: fmt.Sprintf("%s:%d", cfg.RemoteSrvName, cfg.RemoteSrvPort),
	}
	s := server.New(is)
	s.Addr = cfg.SrvAddr
	s.AllowInsecureAuth = lclSrvTls == nil
	s.TLSConfig = lclSrvTls

	is.srv = s
	_imapLog.Info("imap init success at:", cfg.SrvAddr)

	return is, nil
}

func (is *Service) Start() error {
	go func() {
		if is.srv.AllowInsecureAuth {
			_imapLog.Info("imap start success at: ", is.srv.Addr)
			if err := is.srv.ListenAndServe(); err != nil {
				panic(err)
			}
		} else {
			_imapLog.Info("imap with tls start success at:", is.srv.Addr)
			if err := is.srv.ListenAndServeTLS(); err != nil {
				panic(err)
			}
		}
	}()
	return nil
}
