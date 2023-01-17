package proxy_v1

const (
	MaxFindDepth     = 10
	StampSubKey      = "Subject: "
	StampSubSplit    = '\n'
	BlockStampKeyStr = "BlockChain Stamp:"
)

var (
	StampSubSuffix = []byte("======blockchainStamp=====")
)
