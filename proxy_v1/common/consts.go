package common

import "errors"

const (
	CAFileSep        = ";"
	MaxFindDepth     = 10
	StampSubKey      = "Subject: "
	StampSubSplit    = '\n'
	BlockStampKeyStr = "BlockChain Stamp:"
	IMAPSrvName      = "BlockChainStampProtocol"
	IMAPCliVendor    = "StampClient"
)

var (
	TLSErr         = errors.New("no valid tls config")
	StampSubSuffix = []byte("======blockchainStamp=====")
)
