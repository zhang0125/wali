package tron

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/crypto"
)

// Hash goLang sha256 hash algorithm.
func Hash(s []byte) ([]byte, error) {
	h := sha256.New()
	_, err := h.Write(s)
	if err != nil {
		return nil, err
	}
	bs := h.Sum(nil)
	return bs, nil
}

func GetTronAddressFromPriKey(priKey *ecdsa.PrivateKey) []byte {
	publicKey := priKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("get public key fail")
		return nil
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA).String()
	fromAddress = "41" + fromAddress[2:]
	fmt.Println(fromAddress)
	return common.Hex2Bytes(fromAddress)
}
