package go_mail_proxy

import (
	"github.com/blockchainstamp/go-mail-proxy/proxy_v1"
	_ "github.com/btcsuite/btcd/btcutil"
	"sync"
)

var (
	instance ProxyService = nil
	once     sync.Once
)

type ProxyService interface {
	InitByConf(conf any, auth string) error
	Start() error
	ShutDown() error
	StartWithSig(sig chan struct{}) error
}

func Inst() ProxyService {
	once.Do(func() {
		instance = proxy_v1.NewProxy()
	})
	return instance
}
