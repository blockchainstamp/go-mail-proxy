package proxy_v1

import (
	"fmt"
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
	return nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	fmt.Println("Mail from:", from)
	s.env.From = from
	return nil
}

func (s *Session) Rcpt(to string) error {
	fmt.Println("Rcpt to:", to)
	s.env.Tos = append(s.env.Tos, to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	//if b, err := io.ReadAll(r); err != nil {
	//	return err
	//} else {
	//	fmt.Println("Data:", string(b))
	//}
	s.env.Data = r
	return s.sender.SendMail(*s.auth, s.env)
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}
