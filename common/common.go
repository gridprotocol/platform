package common

import (
	"fmt"

	"github.com/grid/contracts/eth"
)

// all contracts addresses
var Contracts = eth.Address{}

// load all contract addresses from json
func init() {
	Contracts = eth.Load("../../grid-contracts/eth/contracts.json")
	fmt.Println("contract addresses:", Contracts)
}
