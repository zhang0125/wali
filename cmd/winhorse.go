package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/zhang0125/wali/tron"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zhang0125/wali/winhorse"
)

const (
	horseTronRpcKey      = "client.rpc"
	horsePriKey          = "client.pri_key"
	horseFeeLimitKey     = "client.fee_limit"
	horseSealContractKey = "horse.address"
	horseContractDataKey = "horse.data"
	horseTestModeKey     = "horse.test"
	horseStartBlockKey   = "horse.start_block"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "horse",
		Short: "Send contract transactions automatically",
		Long:  `Check tron block, send contract transactions when available`,
		Run: func(cmd *cobra.Command, args []string) {
			// get private key
			privateKey, err := crypto.HexToECDSA(viper.GetString(horsePriKey))
			if err != nil {
				fmt.Println("private key error: ", err)
				return
			}
			winNftHorse := &winhorse.WinNFTHorse{
				Done:        false,
				Client:      tron.NewClient(viper.GetString(horseTronRpcKey)),
				PriKey:      privateKey,
				SealAddress: viper.GetString(horseSealContractKey),
				Data:        viper.GetString(horseContractDataKey),
				StartBlock:  viper.GetInt64(horseStartBlockKey),
				FeeLimit:    viper.GetInt64(horseFeeLimitKey),
			}
			if viper.GetBool(horseTestModeKey) {
				winNftHorse.Test()
				return
			}
			//创建监听退出chan
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM,
				syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
			go func() {
				for s := range c {
					switch s {
					case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
						fmt.Println("Start Exit...", s)
						winNftHorse.Stop()
					case syscall.SIGUSR1:
						fmt.Println("usr1 signal", s)
					case syscall.SIGUSR2:
						fmt.Println("usr2 signal", s)
					default:
						fmt.Println("other signal", s)
					}
				}
			}()

			var wg sync.WaitGroup
			wg.Add(1)
			go winNftHorse.Start(&wg)
			wg.Wait()
		},
	})
}
