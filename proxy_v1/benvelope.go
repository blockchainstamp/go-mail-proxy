package proxy_v1

import "io"

type BEnvelope struct {
	From string
	Tos  []string
	Data io.Reader
}

func (env *BEnvelope) WriteTo(w io.Writer) (n int64, err error) {
	return io.Copy(w, env.Data)
}
