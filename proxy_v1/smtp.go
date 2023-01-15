package proxy_v1

type SMTPSrv struct {
}

func NewSMTPSrv(conf *SMTPConf) (*SMTPSrv, error) {
	ss := &SMTPSrv{}
	return ss, nil
}
