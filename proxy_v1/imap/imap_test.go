package imap

import (
	"encoding/json"
	"flag"
	"os"
	"testing"
)

var (
	testImapConf = &Conf{
		RemoteSrvName: "smtp.163.com",
		RemoteSrvPort: 465,
		RemoteSrvCAs:  "rootCAs/163.com.cer;rootCAs/126.com.cer",
		SrvAddr:       ":1143",
		SrvDomain:     "localhost",
	}

	username, password = "", ""
)

func init() {
	flag.StringVar(&username, "usr", "", "")
	flag.StringVar(&password, "pwd", "", "")
}
func TestGenerateIMAPSample(t *testing.T) {
	data, err := json.MarshalIndent(testImapConf, "", "\t")
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile("imap.json.sample", data, 0666); err != nil {
		t.Fatal(err)
	}
}

func TestNewImapServ(t *testing.T) {

	_, err := NewIMAPSrv(testImapConf, nil)
	if err != nil {
		t.Fatal(err)
	}
}
