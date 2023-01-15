package main

import (
	"fmt"
	bp "github.com/blockchainstamp/go-mail-proxy"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version info",
		Long:  `Every software has a version`,
		Run: func(cmd *cobra.Command, args []string) {
			logVersion()
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

func logVersion() {
	fmt.Println("\n==================================================")
	fmt.Printf("Version:\t %s\n", bp.Version)
	fmt.Printf("Build:\t%s\n", bp.BuildTime)
	fmt.Printf("Commit:\t%s\n", bp.Commit)
	fmt.Println("==================================================")
}
