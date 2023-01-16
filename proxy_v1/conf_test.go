package proxy_v1

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestGenerateProxySample(t *testing.T) {
	var c = &Config{
		LogLevel:        uint32(logrus.DebugLevel),
		SMTPConfPath:    "smtp.json",
		BackendConfPath: "backend.json",
	}

	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile("proxy.json.sample", data, 0666); err != nil {
		t.Fatal(err)
	}
}

func TestGenerateSMTPSample(t *testing.T) {
	var c = &SMTPConf{}

	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile("smtp.json.sample", data, 0666); err != nil {
		t.Fatal(err)
	}
}
