package common

import (
	"bufio"
	"bytes"
	"github.com/emersion/go-message/textproto"
	"github.com/sirupsen/logrus"
	"io"
)

var (
	_comLog = logrus.WithFields(logrus.Fields{
		"mode":    "smtp service",
		"package": "common",
	})
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
	_comLog.Debug("===========>msgID:", msgID)
	var headers = map[string][]string{
		BlockStampKeyStr: {"TODO::BlockChain Stamp"},
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
		if depth > MaxFindDepth {
			_comLog.Warn("finding subject exceed max depth:")
			break
		}
		data, err := reader.ReadSlice(StampSubSplit)
		if err != nil {
			_comLog.Debug("finding subject err:", err)
			_, _ = w.Write(data)
			break
		}

		if !bytes.HasPrefix(data, []byte(StampSubKey)) {
			_comLog.Debug("not subject: ", string(data))
			_, _ = w.Write(data)
			continue
		}
		dataLen := len(data)
		if dataLen < 2 {
			_comLog.Warnf("so short[%d] subject!!!", dataLen)
			_, _ = w.Write(data)
			break
		}

		if bytes.Contains(data, StampSubSuffix) {
			_comLog.Warn("no need to add stamp")
			_, _ = w.Write(data)
			break
		}

		var newData []byte
		if data[dataLen-2] == '\r' {
			newData = append(newData, data[:dataLen-2]...)
			newData = append(newData, StampSubSuffix...)
			newData = append(newData, '\r', StampSubSplit)
		} else {
			newData = append(newData, data[:dataLen-1]...)
			newData = append(newData, StampSubSuffix...)
			newData = append(newData, StampSubSplit)
		}
		_comLog.Debug("subject found:", string(newData))

		_, _ = w.Write(newData)
		break
	}

	return io.Copy(w, reader)
}