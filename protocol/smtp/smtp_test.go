package smtp

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/blockchainstamp/go-mail-proxy/utils/common"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	testConf = &Conf{
		SrvAddr:         ":465",
		SrvDomain:       "localhost",
		StampWalletAddr: "BS7UCYp8PSVrjdCn488mz7",
		RemoteConf: map[string]*RemoteConf{
			"163.com": &RemoteConf{
				RemoteCADomain:  "smtp.163.com",
				RemoteSrvName:   "smtp.163.com",
				RemoteSrvPort:   465,
				RemoteSrvCAs:    "rootCAs/163.com.cer",
				ActiveStampAddr: "",
			},
			"126.com": &RemoteConf{
				RemoteCADomain:  "smtp.126.com",
				RemoteSrvName:   "smtp.126.com",
				RemoteSrvPort:   465,
				RemoteSrvCAs:    "rootCAs/126.com.cer",
				ActiveStampAddr: "0x63a39520aE7406d3C63063034eda9EA6d43bFCbb",
			},
			"qq.com": &RemoteConf{
				RemoteCADomain:  "mail.qq.com",
				RemoteSrvName:   "smtp.qq.com",
				RemoteSrvPort:   465,
				ActiveStampAddr: "0xF9Cbfd74808f812a3B8A2337BFC426B2A10Bd74a",
				RemoteSrvCAs:    "rootCAs/qq.com.cer",
			},
		},

		ReadTimeOut:     10,
		WriteTimeOut:    10,
		MaxMessageBytes: 1 << 29,
		MaxRecipients:   50,
	}
	username, password = "", ""
)

func init() {
	flag.StringVar(&username, "usr", "", "")
	flag.StringVar(&password, "pwd", "", "")
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

func testNewSmtpSrv(t *testing.T) *Service {
	ss, err := NewSMTPSrv(testConf, nil)
	if err != nil {
		t.Fatal(err)
	}
	return ss
}

func TestNewSMTPSrv_1(t *testing.T) {
	ss := testNewSmtpSrv(t)
	sig := make(chan struct{}, 1)
	_ = ss.Start(sig)
	<-sig
}

func TestNewSMTPSrv_2(t *testing.T) {
	ss := testNewSmtpSrv(t)
	sig := make(chan struct{}, 1)
	_ = ss.Start(sig)
	<-sig
}

func TestBackendSrv_SendMail(t *testing.T) {
	c := &Conf{
		RemoteConf: map[string]*RemoteConf{
			"ribencong@163.com": &RemoteConf{
				RemoteSrvName: "smtp.163.com",
				RemoteSrvPort: 465,
				RemoteSrvCAs:  "../bin/rootCAs/163.com.cer",
			},
			"ribencong@126.com": &RemoteConf{
				RemoteSrvName: "smtp.126.com",
				RemoteSrvPort: 465,
				RemoteSrvCAs:  "../bin/rootCAs/126.com.cer",
			},
		},
	}

	bs, err := NewSMTPSrv(c, nil)
	if err != nil {
		t.Fatal(err)
	}
	auth := common.Auth{UserName: username, PassWord: password}
	r := strings.NewReader("Subject: Bmail:" + time.Now().String() + "\n\nThis is a test email with blockchain stamp!!!")
	env := &BEnvelope{
		From: username,
		Tos: []string{
			"ribencong@126.com",
		},
		Data: r,
	}
	if err = bs.SendMail(auth, env); err != nil {
		t.Fatal(err)
	}
}
func TestTrimString(t *testing.T) {
	out := strings.TrimLeft(strings.TrimRight("<ctencent_B8CEB2213CA5035DDA981169@qq.com>", ">"), "<")
	fmt.Println(out)
}
