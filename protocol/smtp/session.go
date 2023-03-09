package smtp

import (
	"fmt"
	"github.com/blockchainstamp/go-mail-proxy/utils/common"
	bstamp "github.com/blockchainstamp/go-stamp-wallet"
	"github.com/blockchainstamp/go-stamp-wallet/comm"
	"github.com/emersion/go-smtp"
	"io"
	"strings"
	"time"
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
	err := s.delegate.AUTH(s.auth)
	if err != nil {
		_smtpLog.Infof("user[%s] auth failed:%s", username, err)
		return err
	}
	bstamp.Inst().UpdateStampBalanceAsync(username)
	_smtpLog.Infof("user[%s] auth success:", username)
	return nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	s.env = &BEnvelope{
		From:    from,
		StampID: from,
	}
	stamp := bstamp.Inst().GetStampConf(from)
	if stamp == nil {
		tmp := strings.Split(from, "@")
		if len(tmp) != 2 {
			return fmt.Errorf("invalid email address")
		}
		s.env.StampID = tmp[1]
		stamp = bstamp.Inst().GetStampConf(s.env.StampID)
	}
	if stamp != nil {
		no := 0
		if stamp.IsConsumable {
			no = 1
		}
		s.env.Stamp = &comm.RawStamp{
			SAddr:        comm.StampAddr(stamp.Addr),
			FromMailAddr: from,
			No:           no,
			Time:         time.Now().Unix(),
		}
		_smtpLog.Info("this mail account has stamp")
	}

	_smtpLog.Info("create new envelope from: ", from)
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
