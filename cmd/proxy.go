package main

import (
	"fmt"
	bp "github.com/blockchainstamp/go-mail-proxy"
	"github.com/blockchainstamp/go-mail-proxy/proxy_v1"
	"github.com/blockchainstamp/go-mail-proxy/utils"
	"github.com/blockchainstamp/go-mail-proxy/utils/fdlimit"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

var (
	proxyCmd = &cobra.Command{
		Use:   "proxy",
		Short: "start proxy service",
		Long:  `TODO::`,
		Run:   proxy,
	}

	sigCh      = make(chan os.Signal, 1)
	configPath string
	walletAuth string
)

func init() {
	proxyCmd.Flags().StringVarP(&configPath, "conf",
		"c", "proxy.json", "configure file path --conf||-c [CONFIG_FILE_PATH]")
	proxyCmd.Flags().StringVarP(&walletAuth, "auth",
		"a", "", "--auth||-a [Password Of Current Stamp Wallet]")
	rootCmd.AddCommand(proxyCmd)
}

func initSystem() error {

	if err := os.Setenv("GODEBUG", "netdns=go"); err != nil {
		return err
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(int64(time.Now().Nanosecond()))
	limit, err := fdlimit.Maximum()
	if err != nil {
		return fmt.Errorf("failed to retrieve file descriptor allowance:%s", err)
	}
	_, err = fdlimit.Raise(uint64(limit))
	if err != nil {
		return fmt.Errorf("failed to raise file descriptor allowance:%s", err)
	}
	return nil
}

func waitSignal() {
	pid := strconv.Itoa(os.Getpid())
	fmt.Printf("\n>>>>>>>>>>proxy start at pid(%s)<<<<<<<<<<\n", pid)

	signal.Notify(sigCh,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGUSR1,
		os.Kill,
	)

	//TODO:: process system signal
	for sig := range sigCh {
		fmt.Printf("\n>>>>>>>>>>proxy[%s] finished(%s)<<<<<<<<<<\n", pid, sig)
		return
	}
}

func proxy(cmd *cobra.Command, args []string) {
	if err := initSystem(); err != nil {
		panic(err)
	}
	cnf := &proxy_v1.Config{}
	if err := utils.ReadJsonFile(configPath, cnf); err != nil {
		panic(err)
	}

	if err := bp.Inst().InitByConf(cnf, walletAuth); err != nil {
		panic(err)
	}
	if err := bp.Inst().Start(); err != nil {
		panic(err)
	}
	waitSignal()
	if err := bp.Inst().ShutDown(); err != nil {
		fmt.Println(err)
	}
}
