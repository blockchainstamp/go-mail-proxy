package imap

import (
	"errors"
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
	return u.listMailboxes(subscribed, "*")
}

func (u *User) GetMailbox(name string) (mailbox backend.Mailbox, err error) {
	mailboxes, err := u.listMailboxes(false, name)
	if err != nil {
		return nil, err
	}
	if len(mailboxes) == 0 {
		return nil, errors.New("No such mailbox")
	}

	m := mailboxes[0]
	if err := m.(*Mailbox).ensureSelected(); err != nil {
		return nil, err
	}

	return m, err
}

func (u *User) CreateMailbox(name string) error {
	return u.cli.Create(name)
}

func (u *User) DeleteMailbox(name string) error {
	return u.cli.Delete(name)
}

func (u *User) RenameMailbox(existingName, newName string) error {
	return u.cli.Rename(existingName, newName)
}

func (u *User) Logout() error {
	return u.cli.Logout()
}
