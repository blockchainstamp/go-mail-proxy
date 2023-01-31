package imap

import (
	"crypto/tls"
	"github.com/blockchainstamp/go-mail-proxy/protocol/common"
	"github.com/emersion/go-imap"
	id "github.com/emersion/go-imap-id"
	"github.com/emersion/go-imap/backend"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-imap/server"
	"github.com/sirupsen/logrus"
)

const (
	Version = "1.0.1"
)

var (
	_imapLog = logrus.WithFields(logrus.Fields{
		"mode":    "imap service",
		"package": "imap",
	})
)

type Service struct {
	conf *Conf
	srv  *server.Server
}

func (is *Service) Login(_ *imap.ConnInfo, username, password string) (backend.User, error) {
	conf := is.conf.getRemoteConf(username)
	if conf == nil {
		_imapLog.Warn("no remote tls config for user:", username)
		return nil, common.ConfErr
	}
	u := &User{username: username, password: password}

	c, err := client.DialTLS(conf.remoteSrvAddr, conf.tlsConfig)
	if err != nil {
		_imapLog.Warn("dial failed: ", conf.remoteSrvAddr, err)
		return nil, err
	}

	cli := WrapEXClient(c)
	isID, err := cli.IDCli.SupportID()
	if err != nil {
		return nil, err
	}
	if isID {
		cliID := id.ID{"name": common.IMAPSrvName, "version": Version, "vendor": common.IMAPCliVendor}
		srvID, err := cli.IDCli.ID(cliID)
		if err != nil {
			return nil, err
		}
		tx := srvID["TransID"]
		_imapLog.Debug("imap support id, tx is:", tx)
	}

	if err := cli.Login(username, password); err != nil {
		_imapLog.Warnf("user[%s] login failed:%s", username, err)
		return nil, err
	}
	u.cli = cli
	_imapLog.Infof("user[%s] login success", username)
	return u, nil
}

func NewIMAPSrv(cfg *Conf, lclSrvTls *tls.Config) (*Service, error) {

	if err := cfg.loadRemoteRootCAs(); err != nil {
		return nil, err
	}
	is := &Service{
		conf: cfg,
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
