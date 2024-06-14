package common

import (
	"fmt"

	"github.com/grid/contracts/eth"
)

// all contracts addresses
var Contracts = eth.Address{}

// load all contract addresses from json
func init() {
	Contracts = eth.LoadJSON()
	fmt.Println("contract addresses:", Contracts)
}
