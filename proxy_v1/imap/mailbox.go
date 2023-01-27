package imap

import (
	"time"

	"github.com/emersion/go-imap"
)

type Mailbox struct {
	Subscribed bool
	Messages   []*Message

	name string
	user *User
	info *imap.MailboxInfo
}

func (mbox *Mailbox) ensureSelected() error {
	if mbox.user.cli.Mailbox() != nil && mbox.user.cli.Mailbox().Name == mbox.name {
		return nil
	}

	_, err := mbox.user.cli.Select(mbox.name, false)
	return err
}

func (mbox *Mailbox) Name() string {
	return mbox.name
}

func (mbox *Mailbox) Info() (*imap.MailboxInfo, error) {
	return mbox.info, nil
}

func (mbox *Mailbox) unseenSeqNum() uint32 {
	for i, msg := range mbox.Messages {
		seqNum := uint32(i + 1)

		seen := false
		for _, flag := range msg.Flags {
			if flag == imap.SeenFlag {
				seen = true
				break
			}
		}

		if !seen {
			return seqNum
		}
	}
	return 0
}

func (mbox *Mailbox) Status(items []imap.StatusItem) (*imap.MailboxStatus, error) {
	if mbox.user.cli.Mailbox() != nil && mbox.user.cli.Mailbox().Name == mbox.name {
		mbox := *mbox.user.cli.Mailbox()
		return &mbox, nil
	}

	return mbox.user.cli.Status(mbox.name, items)
}

func (mbox *Mailbox) SetSubscribed(subscribed bool) error {
	if subscribed {
		return mbox.user.cli.Subscribe(mbox.name)
	} else {
		return mbox.user.cli.Unsubscribe(mbox.name)
	}
}

func (mbox *Mailbox) Check() error {
	if err := mbox.ensureSelected(); err != nil {
		return err
	}

	return mbox.user.cli.Check()
}

func (mbox *Mailbox) ListMessages(uid bool, seqSet *imap.SeqSet, items []imap.FetchItem, ch chan<- *imap.Message) error {
	defer close(ch)

	if err := mbox.ensureSelected(); err != nil {
		return err
	}

	messages := make(chan *imap.Message)
	done := make(chan error, 1)
	go func() {
		if uid {
			done <- mbox.user.cli.UidFetch(seqSet, items, messages)
		} else {
			done <- mbox.user.cli.Fetch(seqSet, items, messages)
		}
	}()

	for msg := range messages {
		ch <- msg
	}

	return <-done
}

func (mbox *Mailbox) SearchMessages(uid bool, criteria *imap.SearchCriteria) ([]uint32, error) {
	if err := mbox.ensureSelected(); err != nil {
		return nil, err
	}

	if uid {
		return mbox.user.cli.UidSearch(criteria)
	} else {
		return mbox.user.cli.Search(criteria)
	}
}

func (mbox *Mailbox) CreateMessage(flags []string, date time.Time, body imap.Literal) error {
	defer _imapLog.Debugf("create message with flags%v", flags)
	return mbox.user.cli.Append(mbox.name, flags, date, body)
}

func (mbox *Mailbox) UpdateMessagesFlags(uid bool, seqSet *imap.SeqSet, op imap.FlagsOp, flags []string) error {
	if err := mbox.ensureSelected(); err != nil {
		return err
	}

	flagsInterface := imap.FormatStringList(flags)
	defer _imapLog.Debugf("update message[%s] flags%v uid=%t op=%s", seqSet.String(), flags, uid, op)

	if uid {
		return mbox.user.cli.UidStore(seqSet, imap.StoreItem(op), flagsInterface, nil)
	} else {
		return mbox.user.cli.Store(seqSet, imap.StoreItem(op), flagsInterface, nil)
	}
}

func (mbox *Mailbox) CopyMessages(uid bool, seqSet *imap.SeqSet, destName string) error {
	if err := mbox.ensureSelected(); err != nil {
		return err
	}
	defer _imapLog.Debugf("copy message[%s] to [%s]", seqSet.String(), destName)
	if uid {
		return mbox.user.cli.UidCopy(seqSet, destName)
	} else {
		return mbox.user.cli.Copy(seqSet, destName)
	}
}

func (mbox *Mailbox) Expunge() error {
	if err := mbox.ensureSelected(); err != nil {
		return err
	}
	defer _imapLog.Debugf("expunge from mailbox[%s]", mbox.name)

	return mbox.user.cli.Expunge(nil)
}

func (mbox *Mailbox) MoveMessages(uid bool, seqSet *imap.SeqSet, dest string) error {
	defer _imapLog.Debugf("move message from mailbox[%s] to [%s]", mbox.name, dest)

	if uid {
		return mbox.user.cli.UidMove(seqSet, dest)
	} else {
		return mbox.user.cli.Move(seqSet, dest)
	}
}
