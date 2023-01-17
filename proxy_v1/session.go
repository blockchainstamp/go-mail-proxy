package proxy_v1

import (
	"github.com/emersion/go-smtp"
	"io"
)

type Delegate interface {
	SendMail(auth Auth, env *BEnvelope) error
}

// A Session is returned after EHLO.
type Session struct {
	auth   *Auth
	sender Delegate
	env    *BEnvelope
}

func (s *Session) AuthPlain(username, password string) error {
	s.auth = &Auth{
		UserName: username,
		PassWord: password,
	}
	_proxyLog.Info("session auth for:", username)
	return nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	_proxyLog.Info("Mail from: ", from)
	s.env.From = from
	return nil
}

func (s *Session) Rcpt(to string) error {
	_proxyLog.Info("Rcpt to:", to)
	s.env.Tos = append(s.env.Tos, to)
	return nil
}

func (s *Session) Data(r io.Reader) error {

	//if b, err := io.ReadAll(r); err != nil {
	//	return err
	//} else {
	//	fmt.Println("Data:", string(b))
	//}
	//return nil

	s.env.Data = &BEReader{r}
	return s.sender.SendMail(*s.auth, s.env)
}

func (s *Session) Reset() {
	if s.auth != nil {
		s.env = nil
		_proxyLog.Info("session rest for:", s.auth.UserName)
	}
}

func (s *Session) Logout() error {
	if s.auth != nil {
		_proxyLog.Info("session logout for:", s.auth.UserName)
		s.auth = nil
		s.env = nil
	}
	return nil
}
