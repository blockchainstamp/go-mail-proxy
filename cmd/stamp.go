package main

import (
	"fmt"
	"github.com/blockchainstamp/go-mail-proxy/utils"
	bstamp "github.com/blockchainstamp/go-stamp-wallet"
	"github.com/blockchainstamp/go-stamp-wallet/comm"
	"github.com/spf13/cobra"
)

var (
	stampCmd = &cobra.Command{
		Use:   "stamp",
		Short: "stamp service",
		Long:  `TODO::`,
		Run:   stamp,
	}

	auth       string
	dbPath     string
	walletAddr string
	stampAddr  string
)

func init() {
	stampCmd.Flags().StringVarP(&dbPath, "database",
		"d", "", "--database||-d [STAMP DB PATH]")
	stampCmd.Flags().StringVarP(&auth, "auth",
		"a", "", "--auth|-a [Password Of Stamp Wallet]")

	stampCmd.Flags().StringVar(&walletAddr, "wallet",
		"", "--wallet  [ADDRESS OF Wallet]")
	stampCmd.Flags().StringVar(&stampAddr, "stamp",
		"", "--stamp  [ADDRESS OF Stamp]")

	rootCmd.AddCommand(stampCmd)
}

var stampUsage = `
	Usage:
	bmproxy stamp create-wallet|show
	create-wallet: 
		bmproxy stamp  create-wallet --auth|-a [AUTH]  --database||-d [STAMP DB PATH]
	show: 
		bmproxy stamp show --wallet=[ADDRESS|all] | --stamp=[ADDRESS]
`

func initStamp() bool {
	fi, ok := utils.FileExists(dbPath)
	if !ok || !fi.IsDir() {
		fmt.Println("no such database directory")
		return false
	}

	if err := bstamp.InitSDK(dbPath); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
func stamp(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Print(stampUsage)
		return
	}

	switch args[0] {
	case "create-wallet":
		if len(auth) < 3 {
			if len(auth) == 0 {
				fmt.Println("need --auth|-a [AUTH] ")
				return
			}
			fmt.Println("too short auth --auth")
			return
		}
		if len(dbPath) == 0 {
			fmt.Println("need --database||-d [STAMP DB PATH] ")
			return
		}
		if !initStamp() {
			return
		}
		w, err := bstamp.Inst().CreateWallet(auth)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("---------------------------")
		fmt.Println("Create Success Addr:", w.Address())
		fmt.Println("---------------------------")

	case "show":
		if len(walletAddr) > 0 {
			if !initStamp() {
				return
			}
			if walletAddr == "all" {
				fmt.Println(bstamp.Inst().ListAllWalletAddr())
				return
			}
			w, err := bstamp.Inst().GetWallet(comm.Address(walletAddr))
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(w.Verbose())
			return
		} else if len(stampAddr) > 0 {

		} else {
			fmt.Print(stampUsage)
		}
	default:
		fmt.Print(stampUsage)
		return
	}
}
