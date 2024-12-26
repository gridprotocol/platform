package common

import (
	"github.com/grid/contracts/eth/contracts"
)

// all contracts addresses
var (
	Contracts      contracts.Contracts
	LocalContracts contracts.Local
	SepoContracts  contracts.Sepo
	DevContracts   contracts.Dev
	TestContracts  contracts.Test
)

// load all contract addresses from json
func init() {
	// init for contracts
	Contracts = contracts.Contracts{}
	// init contracts on local chain
	LocalContracts = contracts.Local{}
	LocalContracts.Load()

	// init contracts on sepo chain
	SepoContracts = contracts.Sepo{}
	SepoContracts.Load()

	// init contracts on dev chain
	DevContracts = contracts.Dev{}
	DevContracts.Load()

	// init contracts on test chain
	TestContracts = contracts.Test{}
	TestContracts.Load()
}
