package common

import (
	"crypto/ecdsa"
	"math/big"
	"ground-x/go-gxplatform/common"
	"ground-x/go-gxplatform/core/types"
	"ground-x/go-gxplatform/cmd/istanbul/client"
	"context"
)

var (
	DefaultGasPrice int64 = 0
	DefaultGasLimit int64 = 21000 // the gas of ether tx should be 21000
)

func SendEther(client client.Client, from *ecdsa.PrivateKey, to common.Address, amount *big.Int, nonce uint64) error {
	tx := types.NewTransaction(nonce, to, amount, uint64(DefaultGasLimit) , big.NewInt(DefaultGasPrice), []byte{})
	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, from)
	if err != nil {
		log.Error("Failed to sign transaction", "tx", tx, "err", err)
		return err
	}

	err = client.SendRawTransaction(context.Background(), signedTx)
	if err != nil {
		log.Error("Failed to send transaction", "tx", signedTx, "nonce", nonce, "err", err)
		return err
	}

	return nil
}