package imap

import (
	"errors"
	"github.com/blockchainstamp/go-mail-proxy/protocol/common"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
)

type User struct {
	username  string
	password  string
	cli       *Client
	mailboxes map[string]*Mailbox
}

func (u *User) Username() string {
	return u.username
}
func (u *User) listMailboxes(subscribed bool, name string) ([]backend.Mailbox, error) {
	mailboxes := make(chan *imap.MailboxInfo)
	done := make(chan error, 1)
	go func() {
		if subscribed {
			done <- u.cli.Lsub("", name, mailboxes)
		} else {
			done <- u.cli.List("", name, mailboxes)
		}
	}()

	var list []backend.Mailbox
	for m := range mailboxes {
		list = append(list, &Mailbox{user: u, name: m.Name, info: m})
	}

	return list, <-done
}

func (u *User) ListMailboxes(subscribed bool) (mailboxes []backend.Mailbox, err error) {
	_imapLog.Info("listing all mailbox sub:", subscribed)

	mailboxes, err = u.listMailboxes(subscribed, "*")
	if err != nil {
		return nil, err
	}

	var needCreate = true
	for _, m := range mailboxes {
		if m.Name() == common.StampMailBox {
			_imapLog.Debug("stamp mailbox exist:", m.Name())
			needCreate = false
		}
	}

	if needCreate {
		_imapLog.Info("need to create stamp mailbox")
		if err := u.CreateMailbox(common.StampMailBox); err != nil {
			_imapLog.Warn("failed to create the stamp mailbox:", err)
		}
		return u.listMailboxes(subscribed, "*")
	}
	return mailboxes, nil
}

func (u *User) GetMailbox(name string) (mailbox backend.Mailbox, err error) {
	mailboxes, err := u.listMailboxes(false, name)
	if err != nil {
		_imapLog.Warnf("mailbox[%s] can't be listed:%s", name, err)
		return nil, err
	}
	if len(mailboxes) == 0 {
		_imapLog.Warn("No such mailbox")
		return nil, errors.New("no such mailbox")
	}

	m := mailboxes[0]
	if err := m.(*Mailbox).ensureSelected(); err != nil {
		_imapLog.Warnf("mailbox[%s] can't be selected:%s", name, err)
		return nil, err
	}

	return m, err
}

func (u *User) CreateMailbox(name string) error {
	err := u.cli.Create(name)
	if err != nil {
		_imapLog.Warnf("create mailbox[%s] failed:%s", name, err)
	}
	return err
}

func (u *User) DeleteMailbox(name string) error {
	err := u.cli.Delete(name)
	if err != nil {
		_imapLog.Warnf("delete mailbox[%s] failed:%s", name, err)
	}
	return err
}

func (u *User) RenameMailbox(existingName, newName string) error {
	err := u.cli.Rename(existingName, newName)
	if err != nil {
		_imapLog.Warnf("rename mailbox from[%s] to[%s] failed:%s", existingName, newName, err)
	}
	return err
}

func (u *User) Logout() error {
	_imapLog.Info("user log out:", u.username)
	return u.cli.Logout()
}
