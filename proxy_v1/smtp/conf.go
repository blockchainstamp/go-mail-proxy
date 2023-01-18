package smtp

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/blockchainstamp/go-mail-proxy/proxy_v1/common"
	"os"
	"strings"
)

type Conf struct {
	SrvAddr   string `json:"srv_addr"`
	SrvDomain string `json:"srv_domain"`

	RemoteSrvCAs  string `json:"ca_files"`
	RemoteSrvName string `json:"remote_srv_name"`
	RemoteSrvPort int    `json:"remote_srv_port"`

	MaxMessageBytes int `json:"max_message_bytes"`
	ReadTimeOut     int `json:"read_time_out"`
	WriteTimeOut    int `json:"write_time_out"`
	MaxRecipients   int `json:"max_recipients"`
}

func (sc *Conf) String() string {
	s := "\n=========service============="
	s += "\nServer Addr:\t" + sc.SrvAddr
	s += "\nServer Domain:\t" + sc.SrvDomain
	s += "\nRoot CAs:\t" + sc.RemoteSrvCAs
	s += "\nRemote Server:\t" + sc.RemoteSrvName
	s += fmt.Sprintf("\nRemote Port:\t%d", sc.RemoteSrvPort)
	s += fmt.Sprintf("\nMessage Max:\t%d", sc.MaxMessageBytes)
	s += fmt.Sprintf("\nRead Timout:\t%d", sc.ReadTimeOut)
	s += fmt.Sprintf("\nWrite Timeout:\t%d", sc.WriteTimeOut)
	s += fmt.Sprintf("\nRecipient Max:\t%d", sc.MaxRecipients)
	s += "\n=============================\n"
	return s
}
func (sc *Conf) loadRemoteRootCAs() (*tls.Config, error) {

	fileNames := strings.Split(sc.RemoteSrvCAs, common.CAFileSep)
	if len(fileNames) == 0 {
		return nil, common.TLSErr
	}
	rootCAs := x509.NewCertPool()
	for _, caPath := range fileNames {
		data, err := os.ReadFile(caPath)
		if err != nil {
			return nil, err
		}
		rootCAs.AppendCertsFromPEM(data)
	}
	tlsConf := &tls.Config{
		ServerName: sc.RemoteSrvName,
		RootCAs:    rootCAs,
	}
	return tlsConf, nil
}
