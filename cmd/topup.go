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

// admin topup some credit for an user
var TopupCmd = &cli.Command{
	Name:  "topup",
	Usage: "topup credit",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "address",
			Aliases: []string{"addr"},
			Usage:   "address to topup",
		},
		&cli.StringFlag{
			Name:    "amount",
			Aliases: []string{"a"},
			Usage:   "amount to topup",
		},
	},
	Action: func(ctx *cli.Context) error {
		userAddr := ctx.String("addr")
		amount := ctx.String("amount")
		a, ok := new(big.Int).SetString(amount, 10)
		if !ok {
			return fmt.Errorf("new big int failed")
		}
		creditAddr := comm.Contracts.Credit

		// connect to an eth node with ep
		backend, chainID := eth.ConnETH(eth.Endpoint)
		fmt.Println("chain id:", chainID)

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

		// admin mint credit for user
		tx, err := creditIns.Mint(authAdmin, common.HexToAddress(userAddr), a)
		if err != nil {
			return err
		}

		fmt.Println("waiting for mint tx to be ok")
		// wait tx to complete
		err = eth.CheckTx(eth.Endpoint, tx.Hash(), "")
		if err != nil {
			return err
		}

		return nil
	},
}
