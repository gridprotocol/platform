package cmd

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/grid/contracts/eth"
	"github.com/grid/contracts/go/gtoken"
	comm "github.com/rockiecn/platform/common"
	"github.com/urfave/cli/v2"
)

// admin topup some credit for an user to create orders
var Topup2Cmd = &cli.Command{
	Name:  "topup2",
	Usage: "topup gtoken",
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
			Usage:   "chain to interactivate: local, sepo, dev, test",
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

		// connect to an eth node with ep
		var ep string
		switch chain {
		case "local":
			ep = eth.Ganache
			comm.Contracts = comm.LocalContracts.Contracts
		case "sepo":
			ep = eth.Sepolia
			comm.Contracts = comm.SepoContracts.Contracts
		case "dev":
			ep = eth.DevChain
			comm.Contracts = comm.DevContracts.Contracts
		case "test":
			ep = eth.TestChain
			comm.Contracts = comm.TestContracts.Contracts
		}

		// get credit contract address
		tokenAddr := comm.Contracts.Token

		// connect to chain
		backend, chainID := eth.ConnETH(ep)
		fmt.Println("chain id:", chainID)

		fmt.Println("user addr:", userAddr)
		fmt.Println("gtoken addr:", tokenAddr)

		// get gtoken instance
		gtokenIns, err := gtoken.NewGtoken(common.HexToAddress(tokenAddr), backend)
		if err != nil {
			fmt.Println("new gtoken instance failed:", err)
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
		tx, err := gtokenIns.Transfer(authAdmin, common.HexToAddress(userAddr), v)
		if err != nil {
			return err
		}

		fmt.Println("waiting for transfer tx to be ok: ", tx.Hash())
		// wait tx to complete
		err = eth.CheckTx(ep, tx.Hash(), "")
		if err != nil {
			return err
		}

		fmt.Println("transfer gtoken ok")

		return nil
	},
}
