package smtp

import (
	"github.com/blockchainstamp/go-mail-proxy/protocol/common"
	bstamp "github.com/blockchainstamp/go-stamp-wallet"
	"github.com/blockchainstamp/go-stamp-wallet/comm"
	"github.com/emersion/go-smtp"
	"io"
)

type Delegate interface {
	SendMail(auth common.Auth, env *BEnvelope) error
	AUTH(auth *common.Auth) error
}

type Session struct {
	auth     *common.Auth
	delegate Delegate
	env      *BEnvelope
	conf     *Conf
}

func (s *Session) AuthPlain(username, password string) error {
	s.auth = &common.Auth{
		UserName: username,
		PassWord: password,
	}
	err := s.delegate.AUTH(s.auth)
	if err != nil {
		_smtpLog.Infof("user[%s] auth failed:%s", username, err)
		return err
	}
	_smtpLog.Infof("user[%s] auth success:", username)
	return nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	s.env = &BEnvelope{
		From: from,
	}

	stampAddr := bstamp.Inst().GetActiveStamp(from)
	walletAddr := s.conf.StampWalletAddr

	if len(stampAddr) > 0 && len(walletAddr) > 0 {
		s.env.Stamp = &comm.RawStamp{
			WAddr:        comm.WalletAddr(walletAddr),
			SAdr:         stampAddr,
			FromMailAddr: from,
			No:           1,
		}
		_smtpLog.Info("this mail account has stamp")
	}

	_smtpLog.Info("create new envelope from: ", from)
	bstamp.Inst().UpdateStampBalanceAsync(stampAddr)
	return nil
}

func (s *Session) Rcpt(to string) error {
	_smtpLog.Info("Rcpt to: ", to)
	s.env.Tos = append(s.env.Tos, to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	s.env.Data = r
	return s.delegate.SendMail(*s.auth, s.env)
}

func (s *Session) Reset() {
	if s.auth != nil {
		s.env = nil
		_smtpLog.Info("session rest for:", s.auth.UserName)
	}
}

func (s *Session) Logout() error {
	if s.auth != nil {
		_smtpLog.Info("session logout for:", s.auth.UserName)
		s.auth = nil
		s.env = nil
	}
	return nil
}
