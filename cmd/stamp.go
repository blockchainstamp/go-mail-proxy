package main

import (
	"fmt"
	"github.com/blockchainstamp/go-mail-proxy/utils"
	bstamp "github.com/blockchainstamp/go-stamp-wallet"
	"github.com/blockchainstamp/go-stamp-wallet/comm"
	"github.com/spf13/cobra"
	"os"
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
	walletName string
	stampAddr  string
	export     bool
	walletFile string
)

func init() {
	stampCmd.Flags().StringVarP(&dbPath, "database",
		"d", "", "--database||-d [STAMP DB PATH]")
	stampCmd.Flags().StringVarP(&auth, "auth",
		"a", "", "--auth|-a [Password Of Stamp Wallet]")

	stampCmd.Flags().StringVar(&walletAddr, "wallet",
		"", "--wallet  [ADDRESS OF Wallet]")

	stampCmd.Flags().StringVarP(&walletName, "name",
		"n", "", "--name  [NAME OF Wallet]")

	stampCmd.Flags().StringVar(&stampAddr, "stamp",
		"", "--stamp  [ADDRESS OF Stamp]")
	stampCmd.Flags().BoolVarP(&export, "export", "e", false, "--export|e export wallet data to file")

	stampCmd.Flags().StringVar(&walletFile, "file",
		"", "--file  [FILE OF Wallet]")
	rootCmd.AddCommand(stampCmd)
}

var stampUsage = `
	Usage:
	bmproxy stamp create-wallet|show
	create-wallet: 
		bmproxy stamp  create-wallet --auth|-a [AUTH] --name|-n [NAME]  --database||-d [STAMP DB PATH]
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
	case "import-wallet":
		if len(auth) == 0 {
			fmt.Println("need --auth|-a [AUTH] ")
			return
		}
		if len(walletFile) == 0 {
			fmt.Println("need --file [FILE Of Wallet] ")
			return
		}
		if len(dbPath) == 0 {
			fmt.Println("need --database||-d [STAMP DB PATH] ")
			return
		}
		if !initStamp() {
			return
		}
		bts, err := os.ReadFile(walletFile)
		if err != nil {
			fmt.Println(err)
			return
		}
		w, err := bstamp.Inst().ImportWallet(string(bts), auth)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("import success:", w.Address(), w.EthAddr())
	case "create-wallet":
		if len(auth) < 3 {
			if len(auth) == 0 {
				fmt.Println("need --auth|-a [AUTH] ")
				return
			}
			fmt.Println("too short auth --auth")
			return
		}
		if len(walletName) == 0 {
			fmt.Println("need --name|-n [NAME] ")
			return
		}

		if len(dbPath) == 0 {
			fmt.Println("need --database||-d [STAMP DB PATH] ")
			return
		}
		if !initStamp() {
			return
		}
		w, err := bstamp.Inst().CreateWallet(auth, walletName)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("---------------------------")
		fmt.Println("Create Success:", w.Address(), w.EthAddr())
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
			w, err := bstamp.Inst().GetWallet(comm.WalletAddr(walletAddr))
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(w.Verbose())
			if export {
				_ = os.WriteFile(w.Name()+".export.json", []byte(w.Verbose()), 0666)
			}
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
