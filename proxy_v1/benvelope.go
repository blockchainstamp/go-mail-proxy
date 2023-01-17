package proxy_v1

import (
	"bufio"
	"bytes"
	"io"
)

type BEnvelope struct {
	From string
	Tos  []string
	Data io.Reader
}

func (env *BEnvelope) WriteToWithStamp(w io.Writer) (n int64, err error) {
	_, _ = w.Write([]byte(BlockStampKeyStr + "<wallet address>\n"))
	return io.Copy(w, env.Data)
}

func (env *BEnvelope) WriteTo(w io.Writer) (n int64, err error) {

	var depth = 0
	reader := bufio.NewReader(env.Data)
	for {
		depth++
		if depth > MaxFindDepth {
			_proxyLog.Warn("finding subject exceed max depth:")
			break
		}
		data, err := reader.ReadSlice(StampSubSplit)
		if err != nil {
			_proxyLog.Debug("finding subject err:", err)
			_, _ = w.Write(data)
			break
		}

		if !bytes.HasPrefix(data, []byte(StampSubKey)) {
			_proxyLog.Debug("not subject: ", string(data))
			_, _ = w.Write(data)
			continue
		}
		dataLen := len(data)
		if dataLen < 2 {
			_proxyLog.Warnf("so short[%d] subject!!!", dataLen)
			_, _ = w.Write(data)
			break
		}

		if bytes.Contains(data, StampSubSuffix) {
			_proxyLog.Warn("no need to add stamp")
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
		_proxyLog.Debug("subject found:", string(newData))

		_, _ = w.Write(newData)
		break
	}

	return io.Copy(w, reader)
}
