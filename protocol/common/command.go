package common

const (
	CMDSys   = "SYSTEM"
	CMDProxy = "PROXY"
)

var (
	cmdProcess = make(map[string]CmdHandler)
)

type Command struct {
	Name      string
	Arguments []interface{}
}
type CmdHandler func(cmd Command) any

func RegCmdProc(cmd string, handler CmdHandler) {
	cmdProcess[cmd] = handler
}

func GetCmdProc(cmd string) CmdHandler {
	return cmdProcess[cmd]
}
