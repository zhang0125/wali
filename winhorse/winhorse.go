package winhorse

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"google.golang.org/protobuf/proto"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zhang0125/wali/tron"
)

type WinNFTHorse struct {
	Done        bool
	Client      *tron.Client
	PriKey      *ecdsa.PrivateKey
	SealAddress string
	Data        string
	StartBlock  int64
	FeeLimit    int64
}

func (wf *WinNFTHorse) Test() {
	sealAddressByte := common.Hex2Bytes(wf.SealAddress)
	fromAddress := tron.GetTronAddressFromPriKey(wf.PriKey)

	if fromAddress == nil {
		fmt.Println("cannot get address from private key")
		return
	}
	response, err := wf.Client.TriggerConstantContract(fromAddress, sealAddressByte, common.Hex2Bytes(wf.Data))
	if err != nil {
		fmt.Println("error: ", err)
	} else {
		fmt.Println("response: ", response)
	}
}

func (wf *WinNFTHorse) Start(wg *sync.WaitGroup) {
	for !wf.Done {
		nowBlock, err := wf.Client.GetNowBlock()
		if err != nil {
			panic(err)
		}
		if nowBlock.BlockHeader.RawData.Number+1 >= wf.StartBlock {
			wf.sendTrx()
			wf.Done = true
		}
		fmt.Println("Number:", nowBlock.BlockHeader.RawData.Number, " time:", time.Unix(nowBlock.BlockHeader.RawData.Timestamp/1000, 0))
		time.Sleep(1 * time.Second)
	}

	wg.Done()
}

func (wf *WinNFTHorse) sendTrx() {
	sealAddressByte := common.Hex2Bytes(wf.SealAddress)
	fromAddress := tron.GetTronAddressFromPriKey(wf.PriKey)
	if fromAddress == nil {
		fmt.Println("cannot get address from private key")
		return
	}

	// trigger
	trx, err := wf.Client.TriggerContract(fromAddress, sealAddressByte, common.Hex2Bytes(wf.Data))
	if err != nil {
		fmt.Println("trigger contract error: ", err)
		return
	}
	trx.RawData.FeeLimit = wf.FeeLimit
	rawData, _ := proto.Marshal(trx.GetRawData())
	hash, err := tron.Hash(rawData)
	if err != nil {
		fmt.Println("hash error: ", err)
		return
	}
	signature, err := crypto.Sign(hash, wf.PriKey)
	if err != nil {
		fmt.Println("sign error: ", err)
		return
	}
	trx.Signature = append(trx.GetSignature(), signature)
	err = wf.Client.BroadcastTransaction(context.Background(), trx)
	if err != nil {
		fmt.Println("broadcast error: ", err)
		return
	}
	fmt.Println("broadcast success, txid=:", common.BytesToAddress(hash).String())
}

func (wf *WinNFTHorse) Stop() {
	wf.Done = true
}
