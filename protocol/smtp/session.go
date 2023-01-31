package smtp

import (
	"github.com/blockchainstamp/go-mail-proxy/protocol/common"
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
}

func (s *Session) AuthPlain(username, password string) error {
	s.auth = &common.Auth{
		UserName: username,
		PassWord: password,
	}
	_smtpLog.Info("session auth for:", username)
	return s.delegate.AUTH(s.auth)
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	_smtpLog.Info("Mail from: ", from)
	s.env.From = from
	return nil
}

func (s *Session) Rcpt(to string) error {
	_smtpLog.Info("Rcpt to:", to)
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
