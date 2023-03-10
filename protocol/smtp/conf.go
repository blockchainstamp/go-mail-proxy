package smtp

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/blockchainstamp/go-mail-proxy/utils/common"
	bstamp "github.com/blockchainstamp/go-stamp-wallet"
	"github.com/blockchainstamp/go-stamp-wallet/comm"
	"os"
	"strings"
)

type RemoteConf struct {
	RemoteSrvCAs    string `json:"ca_files"`
	RemoteCADomain  string `json:"ca_domain"`
	AllowNotSecure  bool   `json:"allow_not_secure"`
	RemoteSrvName   string `json:"remote_srv_name"`
	RemoteSrvPort   int    `json:"remote_srv_port"`
	ActiveStampAddr string `json:"active_stamp_addr"`
	tlsConfig       *tls.Config
}

func (rc *RemoteConf) String() string {
	s := "\nRoot CAs:\t" + rc.RemoteSrvCAs
	s += "\nCA Domain:\t" + rc.RemoteCADomain
	s += fmt.Sprintf("\nAllow not security:\t%t", rc.AllowNotSecure)
	s += "\nRemote Server:\t" + rc.RemoteSrvName
	s += "\nStamp Addr:\t" + rc.ActiveStampAddr
	s += fmt.Sprintf("\nRemote Port:\t%d", rc.RemoteSrvPort)
	return s
}

type Conf struct {
	SrvAddr         string                 `json:"srv_addr"`
	SrvDomain       string                 `json:"srv_domain"`
	StampWalletAddr string                 `json:"stamp_wallet_addr"`
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
	s += "\nWallet Addr:\t" + sc.StampWalletAddr
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

func (sc *Conf) prepareAccounts() error {
	for domain, conf := range sc.RemoteConf {
		if err := bstamp.Inst().ConfigStamp(domain, comm.StampAddr(conf.ActiveStampAddr)); err != nil {
			_smtpLog.Error("config stamp failed:", domain, conf.ActiveStampAddr, err)
			return err
		}

		if conf.AllowNotSecure {
			continue
		}
		fileNames := strings.Split(conf.RemoteSrvCAs, common.CAFileSep)
		if len(fileNames) == 0 {
			_smtpLog.Error("no valid ca file:", conf.RemoteSrvCAs)
			return common.TLSErr
		}
		rootCAs := x509.NewCertPool()
		for _, caPath := range fileNames {
			data, err := os.ReadFile(caPath)
			if err != nil {
				_smtpLog.Errorf("read ca file[%s] failed:%s", caPath, err)
				return err
			}
			rootCAs.AppendCertsFromPEM(data)
		}
		tlsConf := &tls.Config{
			ServerName: conf.RemoteCADomain,
			RootCAs:    rootCAs,
		}
		conf.tlsConfig = tlsConf
	}
	return nil
}

func (sc *Conf) getRemoteConf(mailAddr string) *RemoteConf {
	cfg, ok := sc.RemoteConf[mailAddr]
	if ok {
		return cfg
	}
	var addr = strings.Split(mailAddr, common.MailAddrSep)
	if len(addr) != 2 {
		_smtpLog.Warn("invalid email address:", mailAddr)
		return nil
	}
	var domain = addr[1]
	return sc.RemoteConf[domain]
}
