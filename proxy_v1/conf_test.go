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

func TestGenerateBackendSample(t *testing.T) {
	var c = &BackendConf{
		ServerName:  "smtp.gmail.com",
		RootCAFiles: "rootCAs/gmail.com_1.pem;rootCAs/gmail.com_2.pem",
	}

	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile("backend.json.sample", data, 0666); err != nil {
		t.Fatal(err)
	}
}
