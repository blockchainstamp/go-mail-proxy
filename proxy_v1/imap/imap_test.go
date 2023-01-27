package imap

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"github.com/emersion/go-imap"
	id "github.com/emersion/go-imap-id"
	"github.com/emersion/go-imap/client"
	"log"
	"os"
	"testing"
)

var (
	testImapConf = &Conf{
		RemoteSrvName: "imap.163.com",
		RemoteSrvPort: 993,
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

func TestClient(t *testing.T) {
	log.Println("Connecting to server...")
	testRootCAs := x509.NewCertPool()
	testRootCAs.AppendCertsFromPEM(neteaseCert)
	tlsConfig := &tls.Config{
		ServerName: "imap.163.com",
		RootCAs:    testRootCAs,
	}
	// Connect to server
	c, err := client.DialTLS("imap.163.com:993", tlsConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	// Don't forget to logout
	defer c.Logout()
	cli := WrapEXClient(c)
	// Start a TLS session
	isID, err := cli.IDCli.SupportID()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Support ID:", isID)
	if isID {
		cliID := id.ID{"name": "blockchainStamp", "version": "1.0.1", "vendor": "StampClient"}
		srvID, err := cli.IDCli.ID(cliID)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Server ID:", srvID)
	}
	// Now we can login
	if err := cli.Login(username, password); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- cli.List("", "*", mailboxes)
	}()

	log.Println("Mailboxes:")
	for m := range mailboxes {
		log.Println("* " + m.Name)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	// Select INBOX
	mbox, err := cli.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Flags for INBOX:", mbox.Flags)

	// Get the last 4 messages
	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > 3 {
		// We're using unsigned integers here, only substract if the result is > 0
		from = mbox.Messages - 3
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)
	items := []imap.FetchItem{imap.FetchEnvelope}

	messages := make(chan *imap.Message, 10)
	done = make(chan error, 1)
	go func() {
		done <- cli.Fetch(seqset, items, messages)
	}()

	log.Println("Last 4 messages:")
	for msg := range messages {
		log.Println("* " + msg.Envelope.Subject)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	log.Println("Done!")
}

var neteaseCert = []byte(`
-----BEGIN CERTIFICATE-----
MIIGWjCCBUKgAwIBAgIQBkyKaq3fT+CKRKOTe9dYyjANBgkqhkiG9w0BAQsFADBE
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMR4wHAYDVQQDExVH
ZW9UcnVzdCBSU0EgQ04gQ0EgRzIwHhcNMjIwMzI1MDAwMDAwWhcNMjMwNDExMjM1
OTU5WjB1MQswCQYDVQQGEwJDTjERMA8GA1UECBMIemhlamlhbmcxETAPBgNVBAcT
CGhhbmd6aG91MSwwKgYDVQQKEyNOZXRFYXNlIChIYW5nemhvdSkgTmV0d29yayBD
by4sIEx0ZDESMBAGA1UEAwwJKi4xNjMuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOC
AQ8AMIIBCgKCAQEAvWFOPsTGwaR2ppMeFv01BmBViqGrmNf+JF2fMNUq1cC5ZlXF
j617SRrfVqvkhA8YgWeF8El0s1gjHjX2ym+ZWwdF5sH4gX95y05A3IYo6KGyLTvp
04MIiOeT3GJtd9qE0UonjpuHRRgErLugzBdy4UhRMXHNH07z5+OABhWbVIsYJrDQ
awA4T16f/5+ZUmflZ5OP7y5C/Mj3JBa3cIfjdv3XxqqLUa3Dt5IRpIq2o9zaIWqk
cw26JIaa1wuGCiOYAn0j3IqfIFfBwSOTRoc3smrSaaqREPOXOjUhv4InNGNgdRUC
gU/Qj6RFZcisbErKYNUdRhvxOMpNXy5lMJFcsQIDAQABo4IDFTCCAxEwHwYDVR0j
BBgwFoAUJG+RP4mHhw4ywkAY38VM60/ISTIwHQYDVR0OBBYEFHdZ9qFxTou63UPa
wGZXD+PDuo1HMB0GA1UdEQQWMBSCCSouMTYzLmNvbYIHMTYzLmNvbTAOBgNVHQ8B
Af8EBAMCBaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMD0GA1UdHwQ2
MDQwMqAwoC6GLGh0dHA6Ly9jcmwuZGlnaWNlcnQuY24vR2VvVHJ1c3RSU0FDTkNB
RzIuY3JsMD4GA1UdIAQ3MDUwMwYGZ4EMAQICMCkwJwYIKwYBBQUHAgEWG2h0dHA6
Ly93d3cuZGlnaWNlcnQuY29tL0NQUzBxBggrBgEFBQcBAQRlMGMwIwYIKwYBBQUH
MAGGF2h0dHA6Ly9vY3NwLmRpZ2ljZXJ0LmNuMDwGCCsGAQUFBzAChjBodHRwOi8v
Y2FjZXJ0cy5kaWdpY2VydC5jbi9HZW9UcnVzdFJTQUNOQ0FHMi5jcnQwDAYDVR0T
AQH/BAIwADCCAX8GCisGAQQB1nkCBAIEggFvBIIBawFpAHcArfe++nz/EMiLnT2c
Hj4YarRnKV3PsQwkyoWGNOvcgooAAAF/v6HKpQAABAMASDBGAiEA+uqM0ngks7jH
e+UpsC5ZN6AJ8ihjcH3dPQO59tb6FB4CIQD3C5gQiCmAGmogYzsWj//anTfhx2PU
bBSJq3jBcJVQzQB1ADXPGRu/sWxXvw+tTG1Cy7u2JyAmUeo/4SrvqAPDO9ZMAAAB
f7+hyq4AAAQDAEYwRAIgJK6uIk08LUZk1RAAxlqMUcazg6Tg7N1HTsAHWWcdTsUC
IDsv8jIwf+XfNvd/TpHlH7J+aVhii9vih+y/VF8muvRZAHcAs3N3B+GEUPhjhtYF
qdwRCUp5LbFnDAuH3PADDnk2pZoAAAF/v6HK1AAABAMASDBGAiEAv67mHZqyS2ET
TklpVkXdbZHXqgksPrGVP0qgKi/bxDICIQCq8HcjilE3ZhCCG98/GON2wAZf/Fs1
Qec5ebOCnkBysjANBgkqhkiG9w0BAQsFAAOCAQEAIl5PGPoU6+9w61SZN870Znt9
A2Xc4+/UXR0mUVuHeLJPcSpwd4/w1+d1Jnm+dd0vP11bdwvyuc1tekKocdUvwip7
dOCHsqIJmW1T2w0QysBCqLs88Zkrq3HLOGEDVzMk+KZXVUMPgzYgzcUS0rGxjXPi
AyP2XsTEv0rBn8qIwEwcNdgUqPq+XvRs/JZZztWjIQwq3AGmwa/eCSth6PJPqKxU
XQVJLEDLEr/gCwjOdt8HBryBjJ9DDriAhGRQEl5PAmICGhz1/INf/H4gj+tie67X
q4gvvLvTFBhlgbMXwTM19ymgXKeOhXotsBWDCXOlwDR8I1mIamlo0y0nlU3N+Q==
-----END CERTIFICATE-----`)
