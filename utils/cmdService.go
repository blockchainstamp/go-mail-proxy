package utils

import (
	"context"
	"errors"
	"fmt"
	common2 "github.com/blockchainstamp/go-mail-proxy/protocol/common"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"net"
)

type cmdService struct{}

func (s *cmdService) PrintLogLevel(ctx context.Context, request *EmptyRequest) (*CommonResponse, error) {

	return &CommonResponse{
		Msg: logrus.GetLevel().String(),
	}, nil
}

func (s *cmdService) SetLogLevel(ctx context.Context, req *LogLevel) (result *CommonResponse, err error) {
	level, err := logrus.ParseLevel(req.Level)
	if err != nil {
		return nil, err
	}
	logrus.SetLevel(level)
	return &CommonResponse{
		Msg: logrus.GetLevel().String(),
	}, nil
}

func (s *cmdService) ReloadConf(ctx context.Context, request *Config) (*CommonResponse, error) {
	proc := common2.GetCmdProc(common2.CMDProxy)
	if proc == nil {
		return nil, errors.New("no valid processor")
	}

	return &CommonResponse{
		Msg: common2.DefaultCmdSrvAddr,
	}, nil
}

func StartCmdService(addr string) {
	l, err := net.Listen("tcp4", addr)
	if err != nil {
		panic(err)
	}

	cmdServer := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	RegisterCmdServiceServer(cmdServer, &cmdService{})

	reflection.Register(cmdServer)
	fmt.Println("command service start=================>", l.Addr())

	if err := cmdServer.Serve(l); err != nil {
		panic(err)
	}
}

func DialToCmdService(addr string) CmdServiceClient {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := NewCmdServiceClient(conn)

	return client
}
