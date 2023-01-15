package main

import (
	"fmt"
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

	Version   string
	Commit    string
	BuildTime string
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

func logVersion() {
	fmt.Printf("bmproxy %s", Version)
	fmt.Printf("Build Time: %s", BuildTime)
	fmt.Printf("Commit:     %s", Commit)
}
