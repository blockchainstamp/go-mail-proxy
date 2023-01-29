package common

import "errors"

const (
	CAFileSep         = ";"
	MaxFindDepth      = 10
	StampSubKey       = "Subject: "
	StampSubSplit     = '\n'
	BlockStampKey     = "X-Stamp"
	MsgIDKey          = "Message-ID"
	IMAPSrvName       = "BlockChainStampProtocol"
	IMAPCliVendor     = "StampClient"
	StampMailBox      = "StampMailBox" //"区块链邮票"
	INBOXName         = "INBOX"
	SMTPHeaderSize    = 1 << 11
	DefaultCmdSrvAddr = "127.0.0.1:1100"
)

var (
	ConfErr        = errors.New("no config for the user")
	TLSErr         = errors.New("no valid tls config")
	StampSubSuffix = []byte("======blockchainStamp=====")
)
