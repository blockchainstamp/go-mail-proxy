package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "bmproxy",
	Short: "MTA with blockchain mail stamp wallet",
	Long:  `TODO::.`,
	Run:   help,
}

var (
	verbose bool
)

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"print out more debug information")
}
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

var usage = `
	Usage:
`

func help(_ *cobra.Command, _ []string) {
	fmt.Print(usage)
}
