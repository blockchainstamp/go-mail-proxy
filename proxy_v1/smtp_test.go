package proxy_v1

import (
	"encoding/json"
	"github.com/emersion/go-smtp"
	"os"
	"testing"
	"time"
)

var testConf = &SMTPConf{
	Addr:            ":1025",
	Domain:          "localhost",
	ReadTimeOut:     10,
	WriteTimeOut:    10,
	MaxMessageBytes: 1 << 20,
	MaxRecipients:   50,
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

type BE struct {
}

func (p *BE) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &Session{sender: nil, env: &BEnvelope{}}, nil
}
func testNewSmtpSrv(t *testing.T) *SMTPSrv {
	var be = &BE{}
	ss, err := NewSMTPSrv(testConf, be, nil)
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
