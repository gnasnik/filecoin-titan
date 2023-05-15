package contract

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TestConnect(t *testing.T) {
	client, err := ethclient.Dial("https://api.hyperspace.node.glif.io/rpc/v1")
	if err != nil {
		log.Fatal(err)
	}

	address := common.HexToAddress("0x78993e8ed1B4D57De4FC957A096947baB2e14609")
	instance, err := NewContract(address, client)
	if err != nil {
		log.Fatal(err)
	}

	addr := common.HexToAddress("0xeb549F0B9887F4150dbD3bD0A257d99d5E316dBA")

	resp1, err := instance.Store(&bind.TransactOpts{From: addr}, "abc")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp1) // "1.0"

	resp2, err := instance.Get(&bind.CallOpts{From: addr}, "2")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp2) // "1.0"
}

func TestP(t *testing.T) {
	client, err := ethclient.Dial("https://api.hyperspace.node.glif.io/rpc/v1")
	if err != nil {
		log.Println("Dial : ", err)
		return
	}

	address := common.HexToAddress("0x78993e8ed1B4D57De4FC957A096947baB2e14609")
	instance, err := NewContract(address, client)
	if err != nil {
		log.Println("NewContract : ", err)
		return
	}

	privateKey, err := crypto.HexToECDSA("3c3633bfaa3f8cfc2df9169d763eda6a8fb06d632e553f969f9dd2edd64dd11b")
	if err != nil {
		log.Println("HexToECDSA : ", err)
		return
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	// if err != nil {
	// 	log.Println("PendingNonceAt : ", err)
	// 	return
	// }

	// gasPrice, err := client.SuggestGasPrice(context.Background())
	// if err != nil {
	// 	log.Println("SuggestGasPrice : ", err)
	// 	return
	// }

	// auth := bind.NewKeyedTransactor(privateKey)

	signer := types.LatestSignerForChainID(big.NewInt(3141))
	auth := &bind.TransactOpts{
		Signer: func(address common.Address, transaction *types.Transaction) (*types.Transaction, error) {
			return types.SignTx(transaction, signer, privateKey)
		},
		From:    fromAddress,
		Context: context.Background(),
		// GasLimit: gasPrice.Uint64(),
	}

	// auth.Nonce = big.NewInt(int64(nonce))
	// auth.Value = big.NewInt(0)     // in wei
	// auth.GasLimit = uint64(300000) // in units
	// auth.GasPrice = gasPrice

	_, err = instance.Store(auth, "abc123")
	if err != nil {
		log.Println("Store : ", err)
		return
	}

	resp, err := instance.Get(&bind.CallOpts{}, "abc123")
	if err != nil {
		log.Println("Get : ", err)
		return
	}

	log.Println(resp)
}
