package smtp

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/blockchainstamp/go-mail-proxy/utils/common"
	bstamp "github.com/blockchainstamp/go-stamp-wallet"
	"github.com/blockchainstamp/go-stamp-wallet/comm"
	"io"
	"net/textproto"
	"strings"
)

type BEnvelope struct {
	From  string
	Tos   []string
	Stamp comm.StampData
	Data  io.Reader
}

func (env *BEnvelope) WriteTo(w io.Writer) (n int64, err error) {
	var (
		buf  bytes.Buffer
		newH []byte
	)
	tr := io.TeeReader(env.Data, &buf)

	txtR := textproto.NewReader(bufio.NewReaderSize(tr, common.SMTPHeaderSize))
	header, err := txtR.ReadMIMEHeader()
	if err != nil {
		return 0, err
	}

	msgID := header.Get(common.MsgIDKey)
	msgID = strings.TrimLeft(strings.TrimRight(msgID, ">"), "<")

	if env.Stamp != nil {
		env.Stamp.SetMsgID(msgID)
		stamp, err := bstamp.Inst().PostStamp(env.Stamp)
		if err != nil {
			_smtpLog.Warn("sign stamp failed: ", err, msgID)
		} else {
			addH := fmt.Sprintf(common.BlockStampKey+": %s\r\n", stamp.Serial())
			newH = append(newH, []byte(addH)...)
			_smtpLog.Info("append signed stamp success: ", msgID)
		}
	} else {
		_smtpLog.Debug("this mail has no stamp: ", msgID)
	}

	newH = append(newH, buf.Bytes()...)
	_, _ = w.Write(newH)
	return io.Copy(w, env.Data)
}
