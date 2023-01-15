package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var (
	proxyCmd = &cobra.Command{
		Use:   "proxy",
		Short: "start proxy",
		Long:  `TODO::`,
		Run:   proxy,
	}

	sigCh      = make(chan os.Signal, 1)
	configPath string
)

func init() {
	proxyCmd.Flags().StringVarP(&configPath, "conf",
		"c", "proxy.json", "configure file path --conf||-c [CONFIG_FILE_PATH]")
	rootCmd.AddCommand(proxyCmd)
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

	//TODO::

	for sig := range sigCh {
		fmt.Printf("\n>>>>>>>>>>proxy[%s] finished(%s)<<<<<<<<<<\n", pid, sig)
		return
	}
}

func proxy(cmd *cobra.Command, args []string) {
	waitSignal()
}
