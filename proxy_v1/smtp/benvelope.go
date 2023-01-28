package smtp

import (
	"bufio"
	"bytes"
	"github.com/blockchainstamp/go-mail-proxy/proxy_v1/common"
	"github.com/emersion/go-message/textproto"
	"io"
)

type BEnvelope struct {
	From string
	Tos  []string
	Data io.Reader
}

func (env *BEnvelope) WriteTo(w io.Writer) (n int64, err error) {
	tr := io.TeeReader(env.Data, w)
	bufBody := bufio.NewReader(tr)
	subMsgHdr, err := textproto.ReadHeader(bufBody)
	if err != nil {
		return 0, err
	}

	var msgID = subMsgHdr.Get("Message-Id")
	_smtpLog.Debug("msgID:", msgID)
	var headers = map[string][]string{
		common.BlockStampKeyStr: {"TODO::BlockChain Stamp"},
	}
	newH := textproto.HeaderFromMap(headers)
	_ = textproto.WriteHeader(w, newH)
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
