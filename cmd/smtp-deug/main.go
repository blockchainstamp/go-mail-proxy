package main

import (
	"fmt"
	"github.com/emersion/go-smtp"
	"io"
	"io/ioutil"
	"time"
)

type Backend struct{}

func (bkd *Backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &Session{}, nil
}

// A Session is returned after EHLO.
type Session struct{}

func (s *Session) AuthPlain(username, password string) error {
	//if username != "username" || password != "password" {
	//	return errors.New("Invalid username or password")
	//}
	return nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	fmt.Println("Mail from:", from)
	return nil
}

func (s *Session) Rcpt(to string) error {
	fmt.Println("Rcpt to:", to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	if b, err := ioutil.ReadAll(r); err != nil {
		return err
	} else {
		fmt.Println("Data:", string(b))
	}
	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}

func main() {
	be := &Backend{}

	s := smtp.NewServer(be)

	s.Addr = ":1025"
	s.Domain = "localhost"
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	fmt.Println("Starting server at", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}
