package main

import (
	"context"
	"fmt"
	"github.com/blockchainstamp/go-mail-proxy/protocol/common"
	"github.com/blockchainstamp/go-mail-proxy/utils"
	"github.com/spf13/cobra"
)

var (
	cmdCmd = &cobra.Command{
		Use:   "cmd",
		Short: "change system conf",
		Long:  `TODO::`, //TODO::
		Run:   nil,
	}
	logCmd = &cobra.Command{
		Use:   "log",
		Short: "change log setting",
		Long:  `TODO::`, //TODO::
		Run:   logOp,
	}
	confCmd = &cobra.Command{
		Use:   "conf",
		Short: "reload conf",
		Long:  `TODO::`, //TODO::
		Run:   configReload,
	}
	logLevel string
	isShow   bool
	addr     string
	mode     string
)

func init() {
	logCmd.Flags().StringVarP(&logLevel, "set-level", "s", "",
		"log -s|--set-level [panic|fatal|error|warn|info|debug|trace] change current log level")
	logCmd.Flags().BoolVarP(&isShow, "print", "p", false, "-p|--print show current log level")
	cmdCmd.AddCommand(logCmd)

	confCmd.Flags().StringVarP(&mode, "mode", "m", "proxy",
		"conf -m|--mode [proxy|smtp|imap] reload config of MODE")
	confCmd.Flags().BoolVarP(&isShow, "show", "s", false, "-s|--show show current MODE[-m]")

	cmdCmd.AddCommand(confCmd)

	cmdCmd.Flags().StringVarP(&addr, "addr",
		"a", common.DefaultCmdSrvAddr, "conf -a|--addr [Service Network Address]")
	rootCmd.AddCommand(cmdCmd)
}

func configReload(_ *cobra.Command, _ []string) {
	cli := utils.DialToCmdService(addr)
	res, err := cli.ReloadConf(context.Background(), &utils.Config{
		Show: isShow,
		Mode: mode,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res.Msg)
}

func logOp(_ *cobra.Command, _ []string) {
	cli := utils.DialToCmdService(addr)
	if isShow {
		res, err := cli.PrintLogLevel(context.Background(), &utils.EmptyRequest{})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res.Msg)
		return
	}
	if len(logLevel) == 0 {
		fmt.Println("please set log level by -s=[LOG LEVEL]")
		return
	}
	res, err := cli.SetLogLevel(context.Background(), &utils.LogLevel{Level: logLevel})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res.Msg)
}
