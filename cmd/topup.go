package cmd

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/grid/contracts/eth"
	"github.com/grid/contracts/go/credit"
	comm "github.com/rockiecn/platform/common"
	"github.com/urfave/cli/v2"
)

// admin topup some credit for an user to create orders
var TopupCmd = &cli.Command{
	Name:  "topup",
	Usage: "topup credit",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "address",
			Aliases: []string{"a"},
			Usage:   "address to topup",
		},
		&cli.StringFlag{
			Name:    "value",
			Aliases: []string{"v"},
			Usage:   "value to topup",
		},
		&cli.StringFlag{
			Name:    "chain",
			Aliases: []string{"c"},
			Usage:   "chain to interactivate",
			Value:   "local",
		},
	},
	Action: func(ctx *cli.Context) error {
		userAddr := ctx.String("a")
		value := ctx.String("v")
		chain := ctx.String("c")

		// amount to topup
		v, ok := new(big.Int).SetString(value, 10)
		if !ok {
			return fmt.Errorf("new big int failed")
		}
		creditAddr := comm.Contracts.Credit

		// connect to an eth node with ep
		var ep string
		switch chain {
		case "local":
			ep = eth.Endpoint
		case "sepo":
			ep = eth.Endpoint2
		}

		backend, chainID := eth.ConnETH(ep)
		fmt.Println("chain id:", chainID)

		fmt.Println("user addr:", userAddr)
		fmt.Println("credit addr:", creditAddr)

		// get credit instance
		creditIns, err := credit.NewCredit(common.HexToAddress(creditAddr), backend)
		if err != nil {
			fmt.Println("new credit instance failed:", err)
		}

		// make auth to sign and send tx
		authAdmin, err := eth.MakeAuth(chainID, eth.SK0)
		if err != nil {
			return err
		}

		//
		authAdmin.GasLimit = 500000
		// 50 gwei
		authAdmin.GasPrice = new(big.Int).SetUint64(50000000000)

		// admin transfer credit to user
		tx, err := creditIns.Transfer(authAdmin, common.HexToAddress(userAddr), v)
		if err != nil {
			return err
		}

		fmt.Println("waiting for transfer tx to be ok")
		// wait tx to complete
		err = eth.CheckTx(eth.Endpoint, tx.Hash(), "")
		if err != nil {
			return err
		}

		fmt.Println("transfer ok")

		return nil
	},
}
