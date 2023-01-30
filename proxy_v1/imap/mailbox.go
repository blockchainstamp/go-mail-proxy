package imap

import (
	"bufio"
	"bytes"
	"github.com/blockchainstamp/go-mail-proxy/proxy_v1/common"
	"io"
	"net/textproto"
	"time"

	"github.com/emersion/go-imap"
)

type Mailbox struct {
	Subscribed bool
	name       string
	user       *User
	info       *imap.MailboxInfo
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

type bTeeReader struct {
	buf bytes.Buffer
	io.Reader
	io.Writer
}

func (btr *bTeeReader) Write(p []byte) (n int, err error) {
	return btr.buf.Write(p)
}
func (btr *bTeeReader) Read(p []byte) (n int, err error) {
	return btr.buf.Read(p)
}
func (btr *bTeeReader) Len() int {
	return btr.buf.Len()
}

func (btr *bTeeReader) hasStamp() bool {
	txtR := textproto.NewReader(bufio.NewReader(&btr.buf))
	defer btr.buf.Reset()

	headers, err := txtR.ReadMIMEHeader()
	if err != nil {
		_imapLog.Warn("msg header parse err:", err)
		return false
	}
	stamp := headers.Get(common.BlockStampKey)
	if len(stamp) < 4 {
		_imapLog.Info("no stamp found")
		return false
	}
	_imapLog.Debug("stamp found:", stamp)
	return true
}

func (mbox *Mailbox) isStampMail(msg *imap.Message) bool {
	_imapLog.Debug("msg body size:", len(msg.Body))
	for name, literal := range msg.Body {
		if name.BodyPartName.Specifier != imap.HeaderSpecifier {
			continue
		}
		var tr bTeeReader
		_imapLog.Debug("msg header found")
		n, err := io.CopyN(&tr, literal, int64(literal.Len()))
		if err != nil || n < 4 {
			_imapLog.Warn("copy header failed:", err)
			return false
		}
		msg.Body[name] = &tr
		return tr.hasStamp()
	}
	return false
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
	stampSeq := new(imap.SeqSet)
	for msg := range messages {
		if len(msg.Body) > 0 && mbox.name == common.INBOXName {
			if mbox.isStampMail(msg) {
				if uid {
					_imapLog.Debug("stamp mail uid:", msg.Uid)
					stampSeq.AddNum(msg.Uid)
				} else {
					_imapLog.Debug("stamp mail seq:", msg.SeqNum)
					stampSeq.AddNum(msg.SeqNum)
				}
			}
		}
		ch <- msg
	}
	if !stampSeq.Empty() {
		_ = mbox.CopyMessages(uid, stampSeq, common.StampMailBox)
		//err := mbox.MoveMessages(uid, stampSeq, common.StampMailBox)
		_ = mbox.UpdateMessagesFlags(uid, stampSeq, imap.AddFlags, []string{imap.DeletedFlag})
		_ = mbox.Expunge()
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
	defer _imapLog.Debugf("[%s]create message with flags%v", mbox.name, flags)
	return mbox.user.cli.Append(mbox.name, flags, date, body)
}

func (mbox *Mailbox) UpdateMessagesFlags(uid bool, seqSet *imap.SeqSet, op imap.FlagsOp, flags []string) error {
	if err := mbox.ensureSelected(); err != nil {
		return err
	}

	flagsInterface := imap.FormatStringList(flags)
	defer _imapLog.Debugf("[%s]update message[%s] flags%v uid=%t op=%s", mbox.name, seqSet.String(), flags, uid, op)

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
	defer _imapLog.Debugf("[%s]copy message[%s] to [%s]", mbox.name, seqSet.String(), destName)
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
	defer _imapLog.Debugf("move message from mailbox[%s] to [%s] seq:%v", mbox.name, dest, seqSet)

	if uid {
		return mbox.user.cli.UidMove(seqSet, dest)
	} else {
		return mbox.user.cli.Move(seqSet, dest)
	}
}
