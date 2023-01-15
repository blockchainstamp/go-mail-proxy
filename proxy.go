package go_mail_proxy

import (
	"github.com/blockchainstamp/go-mail-proxy/proxy_v1"
	"sync"
)

var (
	instance ProxyService = nil
	once     sync.Once
)

type ProxyService interface {
	InitByConf(confPath string) error
	Start() error
	ShutDown() error
}

func Inst() ProxyService {
	once.Do(func() {
		instance = proxy_v1.NewProxy()
	})
	return instance
}
