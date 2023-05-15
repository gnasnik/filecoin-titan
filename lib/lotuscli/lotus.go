package lotuscli

import (
	"context"
	"fmt"
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
	lotusAddress := "api.node.glif.io"
	authToken := ""

	headers := http.Header{"Authorization": []string{"Bearer " + authToken}}
	url := "http://" + lotusAddress + "/rpc/v0"
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
