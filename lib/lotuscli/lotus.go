package lotuscli

import (
	"context"
	"net/http"

	jsonrpc "github.com/filecoin-project/go-jsonrpc"
	lotusapi "github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
)

// Client lotus client
type Client struct {
	cli *lotusapi.FullNodeStruct
}

// New new a lotus client
func New(lotusAddress, authToken string) (*Client, error) {
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

// ChainHead get head
func (c *Client) ChainHead() (*types.TipSet, error) {
	// Now you can call any API you're interested in.
	return c.cli.ChainHead(context.Background())
}
