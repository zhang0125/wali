package cmd

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/zhang0125/wali/tron/pb"

	"github.com/zhang0125/wali/tron"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	tronRpcKey = "transfer.rpc"

	validatorPriKey = "transfer.pri_key"
	toAddrKey       = "transfer.to_address"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "transfer",
		Short: "Collect validator account balance",
		Run: func(cmd *cobra.Command, args []string) {
			// get private key
			rpc := viper.GetString(tronRpcKey)
			wallet := tron.NewWalletClient(rpc)
			toAddress := viper.GetString(toAddrKey)
			validatorPriKeys := viper.GetStringSlice(validatorPriKey)
			for _, validator := range validatorPriKeys {
				privateKey, err := crypto.HexToECDSA(validator)
				cobra.CheckErr(err)
				account, err := wallet.GetAccount(context.Background(), &pb.Account{
					Address: tron.GetTronAddressFromPriKey(privateKey),
				})
				cobra.CheckErr(err)

				fmt.Println(validator, ":", account.Balance)
				trxEx, err := wallet.CreateTransaction2(context.Background(), &pb.TransferContract{
					OwnerAddress: account.Address,
					ToAddress:    common.Hex2Bytes(toAddress),
					Amount:       account.Balance,
				})
				if err != nil {
					fmt.Println("trigger contract error: ", err)
					return
				}
				signature, err := crypto.Sign(trxEx.Txid, privateKey)
				if err != nil {
					fmt.Println("sign error: ", err)
					return
				}
				trx := trxEx.Transaction
				trx.Signature = append(trx.GetSignature(), signature)
				result, err := wallet.BroadcastTransaction(context.Background(), trx)
				if err != nil {
					fmt.Println("BroadcastTransaction error: ", err)
					return
				}
				if result.Code != pb.Return_SUCCESS {
					fmt.Println("BroadcastTransaction code:", result.Code, " message: ", string(result.Message))
					return
				}
				fmt.Println("broadcast success, txid=:", common.BytesToHash(trxEx.Txid))

			}
		},
	})
}
