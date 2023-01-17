package proxy_v1

const (
	MaxFindDepth  = 10
	StampSubKey   = "Subject: "
	StampSubSplit = '\n'
)

var (
	StampSubSuffix = []byte("======blockchainStamp=====")
)
