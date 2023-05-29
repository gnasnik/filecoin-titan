package recordfile

import (
	"context"
	"fmt"
	"github.com/Filecoin-Titan/titan/lib/lotuscli"
	"testing"
)

func TestGenerateCar(t *testing.T) {

	input := "/home/gnasnik/Pictures/doge.svg"
	outDir := "/home/gnasnik/Pictures/"

	info, err := generateCAR(context.Background(), input, outDir)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(info)
}

/**

Actor ID: 1004
ID Address: t01004
Robust Address: t2bvv546rfigwl6rmze66gwkwgq2gnshd24cyc3sy
Eth Address: 0x40ab732ee16b50314651a1bd2fb04cc0ce9c2d9d
f4 Address: t410ficvxglxbnnidcrsrug6s7mcmydhjylm5nfzdtpy
Return: gxkD7FUCDWveeiVBrL9FmSe8ayrGhozZHHpUQKtzLuFrUDFGUaG9L7BMwM6cLZ0=


*/

func TestStartDeal(t *testing.T) {
	client, err := lotuscli.New()
	if err != nil {
		t.Fatal(err)
	}

	minerID := "t01000"
	from := "t3vikuqj7og7m7rglaf7umremeux5gxmri5mgdfoikawjneiap5pl5ynyjmi4akxv3bmanmb3pnuqnfupywtza"
	input := "/home/gnasnik/Pictures/doge.txt"
	outDir := "/home/gnasnik/Pictures/"
	info, err := generateCAR(context.Background(), input, outDir)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("data cid:", info.DataCid)
	fmt.Println("piece cid:", info.PieceCid)
	fmt.Println("piece size:", info.PieceSize)

	fromAddr := "t410f2krhkdnamf535d7akronkcvwh23w6czmhrseh5a"
	contractAddr := "t01007"

	res, err := InvokeContractByFuncName(context.Background(), client, fromAddr, contractAddr, "addCID(string)", []byte(info.DataCid))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(res))

	res2, err := InvokeContractByFuncName(context.Background(), client, fromAddr, contractAddr, "cidSet(string)", []byte(info.DataCid))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(res2))

	cid, err := createDealProposal(client, minerID, from, info)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(cid)

}
