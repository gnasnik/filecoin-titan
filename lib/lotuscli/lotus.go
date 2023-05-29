package lotuscli

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/filecoin-project/go-address"
	builtintypes "github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/lotus/chain/types/ethtypes"
	"github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"
	"golang.org/x/xerrors"
	"log"
	"net/http"

	jsonrpc "github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/crypto"
	lotusapi "github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
)

// Client lotus client
type Client struct {
	cli *lotusapi.FullNodeStruct
}

// New new a lotus client
func New() (*Client, error) {
	lotusAddress := "127.0.0.1:1234"
	authToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.4GrRGlcdztrzYaTebz-ZyLKP5VugAR4ZUqbtdY7k1nI"

	headers := http.Header{"Authorization": []string{"Bearer " + authToken}}
	url := "http://" + lotusAddress + "/rpc/v1"
	var api lotusapi.FullNodeStruct
	_, err := jsonrpc.NewMergeClient(context.Background(), url, "Filecoin",
		[]interface{}{&api.Internal, &api.CommonStruct.Internal}, headers)
	if err != nil {
		return nil, err
	}

	client := &Client{
		cli: &api,
	}
	return client, nil
}

// StateGetRandomnessFromBeacon get randomness
func (c *Client) StateGetRandomnessFromBeacon() (abi.Randomness, error) {
	ts, err := c.cli.ChainHead(context.Background())
	if err != nil {
		return nil, err
	}

	fmt.Println("height:", ts.Height())
	// c.cli.ClientImport()

	return c.cli.StateGetRandomnessFromBeacon(context.Background(), crypto.DomainSeparationTag_WindowedPoStChallengeSeed, ts.Height(), nil, types.NewTipSetKey())
}

func (c *Client) ImportFile(file lotusapi.FileRef) (*lotusapi.ImportRes, error) {
	return c.cli.ClientImport(context.Background(), file)
}

func (c *Client) StartDealProposal(params *lotusapi.StartDealParams) (*cid.Cid, error) {
	return c.cli.ClientStartDeal(context.Background(), params)
}

func (c *Client) MpoolPushMessage(ctx context.Context, fromAddr, toAddr address.Address, calldata []byte) (ethtypes.EthBytes, error) {
	msg := &types.Message{
		To:       toAddr,
		From:     fromAddr,
		Method:   builtintypes.MethodsEVM.InvokeContract,
		Params:   calldata,
		GasLimit: int64(1000000000),
	}
	log.Println("sending message...")
	smsg, err := c.cli.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return nil, xerrors.Errorf("failed to push message: %w", err)
	}

	log.Println("waiting for message to execute...")
	wait, err := c.cli.StateWaitMsg(ctx, smsg.Cid(), 3, 0, false)
	if err != nil {
		return nil, xerrors.Errorf("error waiting for message: %w", err)
	}

	// check it executed successfully
	if !wait.Receipt.ExitCode.IsSuccess() {
		return nil, xerrors.Errorf("actor execution failed, %+v", wait.Receipt)
	}

	log.Println("Gas used: ", wait.Receipt.GasUsed)
	result, err := cbg.ReadByteArray(bytes.NewBuffer(wait.Receipt.Return), uint64(len(wait.Receipt.Return)))
	if err != nil {
		return nil, xerrors.Errorf("evm result not correctly encoded: %w", err)
	}

	if len(result) > 0 {
		log.Println(hex.EncodeToString(result))
	} else {
		log.Println("OK")
	}

	return nil, nil
}
