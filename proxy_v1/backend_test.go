package proxy_v1

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"mime"
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

	_, err := NewBackendServ(testBackendConf)
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

	bs, err := NewBackendServ(c)
	if err != nil {
		t.Fatal(err)
	}
	auth := Auth{username, password}
	body := bytes.NewBuffer(make([]byte, 1024))
	writeHeader(body, "From", "ribencong@163.com")
	writeHeader(body, "To", "ribencong@126.com")
	env := &BEnvelope{
		From: username,
		Tos: []string{
			"ribencong@126.com",
		},
		Data: strings.NewReader("Subject: Bmail:" + time.Now().String() + "\n\nThis is a test email with blockchain stamp!!!"),
	}
	if err = bs.SendMail(auth, env); err != nil {
		t.Fatal(err)
	}
}
func writeHeader(w io.Writer, k string, v ...string) {
	_, _ = io.WriteString(w, k)
	if len(v) == 0 {
		_, _ = io.WriteString(w, ":\r\n")
		return
	}
	_, _ = io.WriteString(w, ": ")

	// Max header line length is 78 characters in RFC 5322 and 76 characters
	// in RFC 2047. So for the sake of simplicity we use the 76 characters
	// limit.
	charsLeft := 76 - len(k) - len(": ")

	for i, s2 := range v {
		s := mime.BEncoding.Encode(s2, "UTF-8")
		// If the line is already too long, insert a newline right away.
		if charsLeft < 1 {
			if i == 0 {
				_, _ = io.WriteString(w, "\r\n ")
			} else {
				_, _ = io.WriteString(w, ",\r\n ")
			}
			charsLeft = 75
		} else if i != 0 {
			_, _ = io.WriteString(w, ", ")
			charsLeft -= 2
		}

		// While the header content is too long, fold it by inserting a newline.
		for len(s) > charsLeft {
			s = writeLine(w, s, charsLeft)
			charsLeft = 75
		}
		_, _ = io.WriteString(w, s)
		if i := strings.LastIndexByte(s, '\n'); i != -1 {
			charsLeft = 75 - (len(s) - i - 1)
		} else {
			charsLeft -= len(s)
		}
	}
	_, _ = io.WriteString(w, "\r\n")
}

func writeLine(w io.Writer, s string, charsLeft int) string {
	// If there is already a newline before the limit. Write the line.
	if i := strings.IndexByte(s, '\n'); i != -1 && i < charsLeft {
		_, _ = io.WriteString(w, s[:i+1])
		return s[i+1:]
	}

	for i := charsLeft - 1; i >= 0; i-- {
		if s[i] == ' ' {
			_, _ = io.WriteString(w, s[:i])
			_, _ = io.WriteString(w, "\r\n ")
			return s[i+1:]
		}
	}

	// We could not insert a newline cleanly so look for a space or a newline
	// even if it is after the limit.
	for i := 75; i < len(s); i++ {
		if s[i] == ' ' {
			_, _ = io.WriteString(w, s[:i])
			_, _ = io.WriteString(w, "\r\n ")
			return s[i+1:]
		}
		if s[i] == '\n' {
			_, _ = io.WriteString(w, s[:i+1])
			return s[i+1:]
		}
	}

	// Too bad, no space or newline in the whole string. Just write everything.
	_, _ = io.WriteString(w, s)
	return ""
}
