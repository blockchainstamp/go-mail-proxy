package proxy_v1

import (
	"encoding/json"
	"flag"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	testBackendConf = &BackendConf{
		ServerName:  "smtp.163.com",
		ServerPort:  465,
		RootCAFiles: "rootCAs/163.com.cer;rootCAs/126.com.cer",
		ImapAddr:    ":1143",
	}

	username, password = "", ""
)

func init() {
	flag.StringVar(&username, "usr", "", "")
	flag.StringVar(&password, "pwd", "", "")
}
func TestGenerateBackendSample(t *testing.T) {
	data, err := json.MarshalIndent(testBackendConf, "", "\t")
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile("backend.json.sample", data, 0666); err != nil {
		t.Fatal(err)
	}
}

func TestNewBackendServ(t *testing.T) {

	_, err := NewBackendServ(testBackendConf, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackendSrv_SendMail(t *testing.T) {
	c := &BackendConf{
		ServerName:  "smtp.163.com",
		ServerPort:  465,
		RootCAFiles: "../bin/rootCAs/163.com.cer;../bin/rootCAs/126.com.cer",
	}

	bs, err := NewBackendServ(c, nil)
	if err != nil {
		t.Fatal(err)
	}
	auth := Auth{username, password}
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
