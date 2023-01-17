package proxy_v1

import (
	"bufio"
	"bytes"
	"io"
)

type BEReader struct {
	io.Reader
}

func (r *BEReader) appendSuffix(w io.Writer) io.Reader {
	reader := bufio.NewReader(r.Reader)
	for {
		data, err := reader.ReadSlice(StampSubSplit)
		if err != nil {
			_proxyLog.Debug("finding subject err:", err)
			_, _ = w.Write(data)
			break
		}

		if bytes.HasPrefix(data, []byte(StampSubKey)) {
			dataLen := len(data)
			if dataLen > 4 {
				if data[dataLen-2] == '\r' {
					data = append(data[:dataLen-2], StampSubSuffix...)
					data = append(data, '\r', StampSubSplit)
				} else {
					data = append(data[:dataLen-1], StampSubSuffix...)
					data = append(data, StampSubSplit)
				}
			}
			_, _ = w.Write(data)
			_proxyLog.Debug("found subject: ", string(data))
			break
		}
		_proxyLog.Debug("not subject: ", string(data))
		_, _ = w.Write(data)
	}
	return reader
}

type BEnvelope struct {
	From string
	Tos  []string
	Data *BEReader
}

func (env *BEnvelope) WriteTo(w io.Writer) (n int64, err error) {
	r := env.Data.appendSuffix(w)
	return io.Copy(w, r)
}
