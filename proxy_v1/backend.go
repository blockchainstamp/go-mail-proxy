package proxy_v1

import (
	"crypto/tls"
)

type BackendSrv struct {
	tlsCfg *tls.Config
}

func NewBackendServ(conf *BackendConf) (*BackendSrv, error) {
	tlsCfg, err := conf.loadTLSCfg()
	if err != nil {
		return nil, err
	}

	bk := &BackendSrv{tlsCfg: tlsCfg}
	return bk, err
}
