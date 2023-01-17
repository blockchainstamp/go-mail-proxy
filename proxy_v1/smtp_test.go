package proxy_v1

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

var testConf = &SMTPConf{
	Addr:              ":1025",
	Domain:            "localhost",
	ReadTimeOut:       10,
	WriteTimeOut:      10,
	MaxMessageBytes:   1 << 20,
	MaxRecipients:     50,
	AllowInsecureAuth: true,
	TlsKeyPath:        "key.pem",
	TlsCertPath:       "certificate.pem",
}

func TestGenerateSMTPSample(t *testing.T) {

	data, err := json.MarshalIndent(testConf, "", "\t")
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile("smtp.json.sample", data, 0666); err != nil {
		t.Fatal(err)
	}
}

func testNewSmtpSrv(t *testing.T) *SMTPSrv {
	c := &BackendConf{
		ServerName:  "smtp.163.com",
		ServerPort:  465,
		RootCAFiles: "../bin/rootCAs/163.com.cer;../bin/rootCAs/126.com.cer",
	}

	bs, err := NewBackendServ(c)
	if err != nil {
		t.Fatal(err)
	}
	ss, err := NewSMTPSrv(testConf, bs)
	if err != nil {
		t.Fatal(err)
	}
	return ss
}

func TestNewSMTPSrv_1(t *testing.T) {
	ss := testNewSmtpSrv(t)
	_ = ss.Start()
	time.Sleep(30 * time.Second)
}

func TestNewSMTPSrv_2(t *testing.T) {
	ss := testNewSmtpSrv(t)
	_ = ss.Start()
	sig := make(chan bool, 1)
	<-sig
}
