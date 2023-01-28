package common

import "errors"

const (
	CAFileSep     = ";"
	MaxFindDepth  = 10
	StampSubKey   = "Subject: "
	StampSubSplit = '\n'
	BlockStampKey = "X-Stamp"
	IMAPSrvName   = "BlockChainStampProtocol"
	IMAPCliVendor = "StampClient"
	StampMailBox  = "区块链邮票"
	INBOXName     = "INBOX"
)

var (
	ConfErr        = errors.New("no config for the user")
	TLSErr         = errors.New("no valid tls config")
	StampSubSuffix = []byte("======blockchainStamp=====")
)
