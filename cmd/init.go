package main

import (
	"github.com/spf13/cobra"
)

var (
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "init system",
		Long:  `system need a initialization`,
		Run: func(cmd *cobra.Command, args []string) {
			initSys(args)
		},
	}
)

func init() {
	rootCmd.AddCommand(initCmd)
}

func initSys(args []string) {
	if len(args) == 0 {
		return
	}
	switch args[0] {
	case "cert":
		//TODO:: create self ca files key.pem cert.pem
	default:
		return
	}
}
