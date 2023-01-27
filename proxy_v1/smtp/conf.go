package smtp

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/blockchainstamp/go-mail-proxy/proxy_v1/common"
	"os"
	"strings"
)

type RemoteConf struct {
	RemoteSrvCAs  string `json:"ca_files"`
	RemoteSrvName string `json:"remote_srv_name"`
	RemoteSrvPort int    `json:"remote_srv_port"`
	tlsConfig     *tls.Config
}

func (rc *RemoteConf) String() string {
	s := "\nRoot CAs:\t" + rc.RemoteSrvCAs
	s += "\nRemote Server:\t" + rc.RemoteSrvName
	s += fmt.Sprintf("\nRemote Port:\t%d", rc.RemoteSrvPort)
	return s
}

type Conf struct {
	SrvAddr         string                 `json:"srv_addr"`
	SrvDomain       string                 `json:"srv_domain"`
	RemoteConf      map[string]*RemoteConf `json:"remote_conf"`
	MaxMessageBytes int                    `json:"max_message_bytes"`
	ReadTimeOut     int                    `json:"read_time_out"`
	WriteTimeOut    int                    `json:"write_time_out"`
	MaxRecipients   int                    `json:"max_recipients"`
}

func (sc *Conf) String() string {
	s := "\n=========service[smtp]============="
	s += "\nServer Addr:\t" + sc.SrvAddr
	s += "\nServer Domain:\t" + sc.SrvDomain
	s += fmt.Sprintf("\nMessage Max:\t%d", sc.MaxMessageBytes)
	s += fmt.Sprintf("\nRead Timout:\t%d", sc.ReadTimeOut)
	s += fmt.Sprintf("\nWrite Timeout:\t%d", sc.WriteTimeOut)
	s += fmt.Sprintf("\nRecipient Max:\t%d", sc.MaxRecipients)
	for r, conf := range sc.RemoteConf {
		s += fmt.Sprintf("\n---%s---", r)
		s += conf.String()
		s += fmt.Sprintf("\n------------")
	}
	s += "\n=============================\n"
	return s
}

func (sc *Conf) loadRemoteRootCAs() error {
	for _, conf := range sc.RemoteConf {
		fileNames := strings.Split(conf.RemoteSrvCAs, common.CAFileSep)
		if len(fileNames) == 0 {
			return common.TLSErr
		}
		rootCAs := x509.NewCertPool()
		for _, caPath := range fileNames {
			data, err := os.ReadFile(caPath)
			if err != nil {
				return err
			}
			rootCAs.AppendCertsFromPEM(data)
		}
		tlsConf := &tls.Config{
			ServerName: conf.RemoteSrvName,
			RootCAs:    rootCAs,
		}
		conf.tlsConfig = tlsConf
	}
	return nil
}
func (sc *Conf) getRemoteConf(user string) *RemoteConf {
	return sc.RemoteConf[user]
}
