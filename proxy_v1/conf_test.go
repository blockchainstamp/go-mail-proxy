package proxy_v1

import (
	"encoding/json"
	"github.com/blockchainstamp/go-mail-proxy/protocol/common"
	"os"
	"testing"
)

func TestGenerateProxySample(t *testing.T) {
	var c = &Config{
		LogLevel:          "info",
		SMTPConfPath:      "smtp.json",
		IMAPConfPath:      "imap.json",
		CmdSrvAddr:        common.DefaultCmdSrvAddr,
		AllowInsecureAuth: false,
		TlsKeyPath:        "key.pem",
		TlsCertPath:       "certificate.pem",
		StampDBPath:       ".",
	}

	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile("proxy.json.sample", data, 0666); err != nil {
		t.Fatal(err)
	}
}
