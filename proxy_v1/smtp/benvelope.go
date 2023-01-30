package smtp

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/blockchainstamp/go-mail-proxy/proxy_v1/common"
	bstamp "github.com/blockchainstamp/go-stamp-wallet"
	"io"
	"net/textproto"
)

type BEnvelope struct {
	From string
	Tos  []string
	Data io.Reader
}

func (env *BEnvelope) WriteTo(w io.Writer) (n int64, err error) {

	var buf bytes.Buffer
	tr := io.TeeReader(env.Data, &buf)

	txtR := textproto.NewReader(bufio.NewReaderSize(tr, common.SMTPHeaderSize))
	header, err := txtR.ReadMIMEHeader()
	if err != nil {
		return 0, err
	}
	msgID := header.Get(common.MsgIDKey)
	_smtpLog.Debug("msgID:", msgID)
	var newH []byte
	stamp, err := bstamp.Inst().CreateStamp(env.From, msgID)
	if err != nil {
		_smtpLog.Warn("create stamp failed:", err)
	} else {
		addH := fmt.Sprintf(common.BlockStampKey+": %s\r\n", stamp.Serial())
		newH = append(newH, []byte(addH)...)
	}
	newH = append(newH, buf.Bytes()...)
	_, _ = w.Write(newH)
	return io.Copy(w, env.Data)
}

func (env *BEnvelope) WriteToOld(w io.Writer) (n int64, err error) {

	var depth = 0
	reader := bufio.NewReader(env.Data)
	for {
		depth++
		if depth > common.MaxFindDepth {
			_smtpLog.Warn("finding subject exceed max depth:")
			break
		}
		data, err := reader.ReadSlice(common.StampSubSplit)
		if err != nil {
			_smtpLog.Debug("finding subject err:", err)
			_, _ = w.Write(data)
			break
		}

		if !bytes.HasPrefix(data, []byte(common.StampSubKey)) {
			_smtpLog.Debug("not subject: ", string(data))
			_, _ = w.Write(data)
			continue
		}
		dataLen := len(data)
		if dataLen < 2 {
			_smtpLog.Warnf("so short[%d] subject!!!", dataLen)
			_, _ = w.Write(data)
			break
		}

		if bytes.Contains(data, common.StampSubSuffix) {
			_smtpLog.Warn("no need to add stamp")
			_, _ = w.Write(data)
			break
		}

		var newData []byte
		if data[dataLen-2] == '\r' {
			newData = append(newData, data[:dataLen-2]...)
			newData = append(newData, common.StampSubSuffix...)
			newData = append(newData, '\r', common.StampSubSplit)
		} else {
			newData = append(newData, data[:dataLen-1]...)
			newData = append(newData, common.StampSubSuffix...)
			newData = append(newData, common.StampSubSplit)
		}
		_smtpLog.Debug("subject found:", string(newData))

		_, _ = w.Write(newData)
		break
	}

	return io.Copy(w, reader)
}
