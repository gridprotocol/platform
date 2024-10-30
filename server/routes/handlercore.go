package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/grid/contracts/eth"
	"github.com/grid/contracts/go/credit"
	"github.com/grid/contracts/go/market"
	"github.com/grid/contracts/go/registry"
	"github.com/grid/contracts/go/version"
	comm "github.com/rockiecn/platform/common"
	"github.com/rockiecn/platform/lib/kv"
	"github.com/rockiecn/platform/lib/utils"
)

type HandlerCore struct {
	LocalDB *kv.Database
}

// pay info when recharge credit, need to be stored in db
type PayInfo struct {
	PIKey  string `json:"PayInfoKey"`
	TIKey  string `json:"TransferInfoKey"`
	Owner  string `json:"Owner"`
	Credit int64  `json:"Credit"`
	TxHash string `json:"TxHash"`
}

// info about a transfer
type TransferInfo struct {
	TIKey       string `json:"TxInfoKey"`
	TxHash      string `json:"TxHash"`
	From        string `json:"From"`
	To          string `json:"To"`
	Value       string `json:"Value"`
	TxConfirmed bool   `json:"TxConfirmed"`
	CreditSaved bool   `json:"CreditSaved"`
}

// ShowAccount godoc
//
//	@Summary		welcome
//	@Description	welcome api
//	@Tags			Welcome
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	string	"file id"
//	@Router			/ [get]
func (hc *HandlerCore) RootHandler(c *gin.Context) {
	c.String(http.StatusOK, "Welcome Server")
}

// handler of cp login
/*
//	@Summary		Regist CP
//	@Description	Regist CP
//	@Tags			RegistCP
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			name		formData	string	true	"name"
//	@Param			address		formData	string	true	"address"
//	@Param			endpoint	formData	string	true	"endpoint"
//	@Param			numCPU		formData	string	true	"num cpu"
//	@Param			priCPU		formData	string	true	"price cpu"
//	@Param			numGPU		formData	string	true	"num"
//	@Param			priGPU		formData	string	false	"price"
//	@Param			numDisk		formData	string	true	"num"
//	@Param			priDisk		formData	string	false	"price"
//	@Param			numMem		formData	string	true	"num"
//	@Param			priMem		formData	string	false	"price"
//	@Success		200			{object}	string	"regist OK"
//	@Failure		400			{object}	string	"bad request"
//	@Router			/registcp [post]
*/
func (hc *HandlerCore) RegistCPHandler(c *gin.Context) {
	// provider wallet address
	cpAddr := c.PostForm("address")

	// get signeTx from input
	tx := c.PostForm("tx")

	//log.Println("tx: ", tx)

	// transfer to types.Transaction
	signedTx := new(types.Transaction)
	signedTx.UnmarshalJSON([]byte(tx))
	log.Println("signed tx: ", signedTx)

	// connect to an eth client
	log.Println("connecting client")
	client, err := ethclient.Dial(Chain_Endpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// send tx to register a cp
	log.Println("sending tx")
	// send a tx to client
	if err := client.SendTransaction(context.Background(), signedTx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"send error": err.Error()})
		return
	}

	// wait tx ok
	logger.Info("waiting for tx to be ok")
	eth.CheckTx(Chain_Endpoint, signedTx.Hash(), "")

	// get cp's reg info

	// get registry instance
	regIns, err := registry.NewRegistry(common.HexToAddress(comm.Contracts.Registry), client)
	if err != nil {
		log.Println("new registry instance failed: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// check cp's reg info
	regInfo, err := regIns.Get(&bind.CallOpts{}, common.HexToAddress(cpAddr))
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("regInfo:", regInfo)

	// response to client
	c.JSON(http.StatusOK, gin.H{"response": "regist OK"})

}

// string memory ip,
// string memory domain,
// string memory port,
// uint64 cpu,
// uint64 gpu,
// uint64 mem,
// uint64 disk,
// uint64 pcpu,
// uint64 pgpu,
// uint64 pmem,
// uint64 pdisk

// revise cp info
func (hc *HandlerCore) ReviseHandler(c *gin.Context) {
	// provider wallet address
	cpAddr := c.PostForm("address")

	// get signeTx from input
	tx := c.PostForm("tx")

	//log.Println("tx: ", tx)

	// transfer to types.Transaction
	signedTx := new(types.Transaction)
	signedTx.UnmarshalJSON([]byte(tx))
	log.Println("signed tx: ", signedTx)

	// connect to an eth client
	log.Println("connecting client")
	client, err := ethclient.Dial(Chain_Endpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// send tx to register a cp
	log.Println("sending tx")
	// send a tx to client
	if err := client.SendTransaction(context.Background(), signedTx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// wait tx ok
	logger.Info("waiting for set to be ok")
	eth.CheckTx(Chain_Endpoint, signedTx.Hash(), "")

	// get cp's reg info

	// get registry instance
	regIns, err := registry.NewRegistry(common.HexToAddress(comm.Contracts.Registry), client)
	if err != nil {
		log.Println("new registry instance failed: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// check cp's reg info
	regInfo, err := regIns.Get(&bind.CallOpts{}, common.HexToAddress(cpAddr))
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("regInfo:", regInfo)

	// response to client
	c.JSON(http.StatusOK, gin.H{"response": "revise OK"})

}

// handler for list cp nodes
// ListCPHandler godoc
//
//	@Summary		List all providers
//	@Description	list all providers
//	@Tags			Listcps
//	@Accept			json
//	@Produce		json
//	@Param			start	query		string	true	"start"
//	@Param			num		query		string	true	"number"
//	@Success		200		{object}	[]CPInfo
//	@Failure		404		{object}	string	"page not found"
//	@Router			/listcp/ [get]
func (hc *HandlerCore) ListCPHandler(c *gin.Context) {
	start := c.Query("start")
	num := c.Query("num")

	s, _ := utils.StringToUint64(start)
	n, _ := utils.StringToUint64(num)

	// connect to an eth node with ep
	logger.Info("connecting chain")
	backend, chainID := eth.ConnETH(Chain_Endpoint)
	logger.Info("chain id:", chainID)

	logger.Info("registry address: ", comm.Contracts.Registry)
	// get contract instance
	registryIns, err := registry.NewRegistry(common.HexToAddress(comm.Contracts.Registry), backend)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("new registry instance failed: %s", err.Error()).Error()})
		return
	}

	// get cp list
	list, _, err := registryIns.ListCP(&bind.CallOpts{}, s, n)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("get cp keys failed: %s", err.Error()).Error()})
		return
	}
	//logger.Info("cp keys:", keys)

	// response key list
	c.JSON(http.StatusOK, list)
}

// handler for get a cp node
// GetCPHandler godoc
//
//	@Summary		get a provider
//	@Description	get a provider's info
//	@Tags			Get Provider Info
//	@Accept			json
//	@Produce		json
//	@Param			address	query		string	true	"address of a provider"
//	@Success		200		{object}	CPInfo
//	@Failure		404		{object}	string	"page not found"
//	@Failure		500		{object}	string	"internal server error"
//
//	@Router			/getcp/ [get]
func (hc *HandlerCore) GetCPHandler(c *gin.Context) {

	// provider address from param
	cpaddr := c.Query("address")

	// connect to an eth node with ep
	logger.Info("connecting chain")
	backend, chainID := eth.ConnETH(Chain_Endpoint)
	logger.Info("chain id:", chainID)

	logger.Debug("cp addr:", cpaddr)
	logger.Debug("registry addr:", comm.Contracts.Registry)

	// get contract instance
	regIns, err := registry.NewRegistry(common.HexToAddress(comm.Contracts.Registry), backend)

	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("new token instance failed: %v", err))
		return
	}

	// get balance of addr2
	regInfo, err := regIns.Get(&bind.CallOpts{}, common.HexToAddress(cpaddr))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	logger.Info("registry info:", regInfo)

	// response node list
	c.JSON(http.StatusOK, regInfo)
}

// user approve the order fee to market before create order
func (hc *HandlerCore) ApproveHandler(c *gin.Context) {
	// tx in form
	txData := c.PostForm("tx")
	log.Println("tx: ", txData)

	// transfer to types.Transaction
	signedTx := new(types.Transaction)
	signedTx.UnmarshalJSON([]byte(txData))
	log.Println("signed tx: ", signedTx)

	// connect to an eth client
	log.Println("connecting client")
	client, err := ethclient.Dial(Chain_Endpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// send approve tx

	log.Println("sending tx")
	// send a tx to client
	if err := client.SendTransaction(context.Background(), signedTx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// wait tx ok
	logger.Info("waiting for set to be ok")
	eth.CheckTx(Chain_Endpoint, signedTx.Hash(), "")

	// response
	c.JSON(http.StatusOK, gin.H{
		"msg": "user approve ok",
	})
}

// check credit allowance
// AllowanceHandler godoc
//
//	@Summary		Check Allowance
//	@Description	check the allowance between an owner and a spender
//	@Tags			Allowance
//	@Accept			json
//	@Produce		json
//	@Param			owner	query		string	true	"owner"
//	@Param			spender	query		string	true	"spender"
//	@Success		200		{object}	int
//	@Failure		404		{object}	string	"page not found"
//	@Router			/allowance/ [get]
func (hc *HandlerCore) AllowanceHandler(c *gin.Context) {
	// param in form
	owner := c.Query("owner")
	spender := c.Query("spender")

	// connect to an eth node with ep
	logger.Info("connecting chain")
	backend, chainID := eth.ConnETH(Chain_Endpoint)
	logger.Info("chain id:", chainID)

	// get contract instance
	creditIns, err := credit.NewCredit(common.HexToAddress(comm.Contracts.Credit), backend)

	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("new credit instance failed: %v", err))
		return
	}

	// get allowance
	allow, err := creditIns.Allowance(&bind.CallOpts{}, common.HexToAddress(owner), common.HexToAddress(spender))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logger.Infof("allowance:", allow)

	// response
	c.JSON(http.StatusOK, gin.H{"allowance": allow})
}

// handler of create order
/*
//	@Summary		Create order
//	@Description	create an order
//	@Tags			CreateOrder
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			userAddress	formData	string	true	"user address"
//	@Param			cpAddress	formData	string	true	"cpAddress"
//	@Param			numCPU		formData	string	true	"num cpu"
//	@Param			priCPU		formData	string	true	"price cpu"
//	@Param			numGPU		formData	string	true	"num"
//	@Param			priGPU		formData	string	false	"price"
//	@Param			numDisk		formData	string	true	"num"
//	@Param			priDisk		formData	string	false	"price"
//	@Param			numMem		formData	string	true	"num"
//	@Param			priMem		formData	string	false	"price"
//	@Param			duration	formData	string	true	"duration"
//	@Success		200			{object}	string	"regist OK"
//	@Failure		400			{object}	string	"bad request"
//	@Router			/createorder [post]
*/
func (hc *HandlerCore) CreateOrderHandler(c *gin.Context) {
	// tx data in form
	txData := c.PostForm("tx")

	// transfer to types.Transaction
	signedTx := new(types.Transaction)
	signedTx.UnmarshalJSON([]byte(txData))
	log.Println("signed tx: ", signedTx)

	// connect to an eth client
	log.Println("connecting client")
	client, err := ethclient.Dial(Chain_Endpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// send createorder tx

	log.Println("sending tx")
	// send a tx to client
	if err := client.SendTransaction(context.Background(), signedTx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// wait tx ok
	logger.Info("waiting for set to be ok")
	eth.CheckTx(Chain_Endpoint, signedTx.Hash(), "")

	// response
	c.JSON(http.StatusOK, gin.H{
		"msg": "create order ok",
	})
}

func (hc *HandlerCore) UserConfirmHandler(c *gin.Context) {
	// tx data in form
	txData := c.PostForm("tx")

	// transfer to types.Transaction
	signedTx := new(types.Transaction)
	signedTx.UnmarshalJSON([]byte(txData))
	log.Println("signed tx: ", signedTx)

	// connect to an eth client
	log.Println("connecting client")
	client, err := ethclient.Dial(Chain_Endpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// send createorder tx

	log.Println("sending tx")
	// send a tx to client
	if err := client.SendTransaction(context.Background(), signedTx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// wait tx ok
	logger.Info("waiting for set to be ok")
	eth.CheckTx(Chain_Endpoint, signedTx.Hash(), "")

	// response
	c.JSON(http.StatusOK, gin.H{
		"msg": "user confirm ok",
	})
}

func (hc *HandlerCore) UserCancelHandler(c *gin.Context) {
	// tx data in form
	txData := c.PostForm("tx")

	// transfer to types.Transaction
	signedTx := new(types.Transaction)
	signedTx.UnmarshalJSON([]byte(txData))
	log.Println("signed tx: ", signedTx)

	// connect to an eth client
	log.Println("connecting client")
	client, err := ethclient.Dial(Chain_Endpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("sending tx")
	// send a tx to client
	if err := client.SendTransaction(context.Background(), signedTx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// wait tx ok
	logger.Info("waiting for set to be ok")
	eth.CheckTx(Chain_Endpoint, signedTx.Hash(), "")

	// response
	c.JSON(http.StatusOK, gin.H{
		"msg": "user cancel ok",
	})
}

// handler for getOrder
// GetOrderHandler godoc
//
//	@Summary		Get Order
//	@Description	get an order info
//	@Tags			Get Order
//	@Accept			json
//	@Produce		json
//	@Param			user	query		string	true	"user"
//	@Param			cp		query		string	true	"cp"
//	@Success		200		{object}	int
//	@Failure		404		{object}	string	"page not found"
//	@Router			/getorder/ [get]
func (hc *HandlerCore) GetOrderHandler(c *gin.Context) {

	// user and cp
	userAddr := c.Query("user")
	cpAddr := c.Query("cp")

	logger.Debug("user:", userAddr)
	logger.Debug("cp:", cpAddr)

	// connect to an eth node with ep
	logger.Info("connecting chain")
	backend, chainID := eth.ConnETH(Chain_Endpoint)
	logger.Info("chain id:", chainID)

	logger.Debug("market:", comm.Contracts.Market)

	// get contract instance
	marketIns, err := market.NewMarket(common.HexToAddress(comm.Contracts.Market), backend)

	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("new market instance failed: %v", err))
		return
	}

	// get order with user and cp
	orderInfo, err := marketIns.GetOrder(&bind.CallOpts{}, common.HexToAddress(userAddr), common.HexToAddress(cpAddr))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logger.Info("order info:", orderInfo)

	// response node list
	c.JSON(http.StatusOK, orderInfo)
}

/*
// value of an order
func (hc *HandlerCore) ValueOrderHandler(c *gin.Context) {

	// user and cp
	userAddr := c.Query("user")
	cpAddr := c.Query("cp")

	logger.Debug("user:", userAddr)
	logger.Debug("cp:", cpAddr)

	// connect to an eth node with ep
	logger.Info("connecting chain")
	backend, chainID := eth.ConnETH(Chain_Endpoint)
	logger.Info("chain id:", chainID)

	logger.Debug("market:", comm.Contracts.Market)

	// get contract instance
	marketIns, err := market.NewMarket(common.HexToAddress(comm.Contracts.Market), backend)

	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("new market instance failed: %v", err))
		return
	}

	// get value
	value, err := marketIns.ValueOrder(&bind.CallOpts{}, common.HexToAddress(userAddr), common.HexToAddress(cpAddr))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logger.Info("order value:", value)

	// response node list
	c.JSON(http.StatusOK, value)
}
*/

// handler for get order list for user or cp
//
//	@Summary		List order
//	@Description	list an order
//	@Tags			ListOrder
//	@Accept			json
//	@Produce		json
//	@Param			role	query		string	true	"user or provider"
//	@Param			address	query		string	true	"address"
//	@Success		200		{object}	string	"list OK"
//	@Failure		400		{object}	string	"bad request"
//	@Router			/listorder [get]
// func (hc *HandlerCore) ListOrderHandler(c *gin.Context) {

// 	// connect to an eth node with ep
// 	logger.Info("connecting chain")
// 	backend, chainID := eth.ConnETH(Chain_Endpoint)
// 	logger.Info("chain id:", chainID)

// 	// get contract instance
// 	marketIns, err := market.NewMarket(common.HexToAddress(comm.Contracts.Market), backend)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("new market instance failed: %s", err.Error())})
// 	}

// 	// get key list
// 	keys, err := marketIns.GetKeys(&bind.CallOpts{})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("get order keys failed: %s", err.Error())})
// 	}
// 	//logger.Info("cp keys:", keys)

// 	// response key list
// 	c.JSON(http.StatusOK, keys)
// }

// handler for get provider list for an user
//
//	@Summary		Get Pro List
//	@Description	get the provider list of an user
//	@Tags			GetProList
//	@Accept			json
//	@Produce		json
//	@Param			user	query		string	true	"user address"
//	@Success		200		{object}	string	"list OK"
//	@Failure		400		{object}	string	"bad request"
//	@Router			/listorder [get]
func (hc *HandlerCore) GetCPSHandler(c *gin.Context) {
	userAddr := c.Query("user")

	// connect to an eth node with ep
	logger.Info("connecting chain")
	backend, chainID := eth.ConnETH(Chain_Endpoint)
	logger.Info("chain id:", chainID)

	// get contract instance
	marketIns, err := market.NewMarket(common.HexToAddress(comm.Contracts.Market), backend)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("new market instance failed: %s", err.Error())})
	}

	// get key list
	keys, err := marketIns.GetProList(&bind.CallOpts{}, common.HexToAddress(userAddr))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("get order keys failed: %s", err.Error())})
		return
	}
	//logger.Info("cp keys:", keys)

	// response key list
	c.JSON(http.StatusOK, keys)
}

// get a node info
// GetNodeHandler godoc
//
//	@Summary		Node
//	@Description	Get a node of a cp with node id
//	@Tags			Get Node
//	@Accept			json
//	@Produce		json
//	@Param			cp	query		string	true	"cp address"
//	@Param			id	query		string	true	"node id"
//	@Success		200	{object}	int
//	@Failure		404	{object}	string	"page not found"
//	@Router			/node/ [get]
func (hc *HandlerCore) GetNodeHandler(c *gin.Context) {
	cp := c.Query("cp")
	id := c.Query("id")

	// connect to an eth node with ep
	logger.Info("connecting chain")
	backend, chainID := eth.ConnETH(Chain_Endpoint)
	logger.Info("chain id:", chainID)

	logger.Info("get registry instance")
	// get contract instance
	regIns, err := registry.NewRegistry(common.HexToAddress(comm.Contracts.Registry), backend)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("new market instance failed: %s", err.Error())})
		return
	}

	logger.Info("call list node")

	// string to big
	id64, _ := utils.StringToUint64(id)

	// get node info
	node, err := regIns.GetNode(&bind.CallOpts{}, common.HexToAddress(cp), id64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("get node failed: %s", err.Error())})
		return
	}

	// response
	c.JSON(http.StatusOK, node)
}

// list node for a proivder
// GetNodesHandler godoc
//
//	@Summary		Nodes
//	@Description	Get all nodes of this provider
//	@Tags			Get Nodes
//	@Accept			jsonb
//	@Produce		json
//	@Param			cp	query		string	true	"cp address"
//	@Success		200	{object}	int
//	@Failure		404	{object}	string	"page not found"
//	@Router			/nodes/ [get]
func (hc *HandlerCore) GetNodesHandler(c *gin.Context) {
	cp := c.Query("cp")

	// connect to an eth node with ep
	logger.Info("connecting chain")
	backend, chainID := eth.ConnETH(Chain_Endpoint)
	logger.Info("chain id:", chainID)

	logger.Info("get registry instance")
	// get contract instance
	regIns, err := registry.NewRegistry(common.HexToAddress(comm.Contracts.Registry), backend)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("new market instance failed: %s", err.Error())})
		return
	}

	logger.Info("call list node")

	// get node list
	nodes, err := regIns.ListNode(&bind.CallOpts{}, common.HexToAddress(cp))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("list node failed: %s", err.Error())})
		return
	}
	//logger.Info("cp keys:", keys)

	// response key list
	c.JSON(http.StatusOK, nodes)
}

// get orders for an user
// GetOrdersHandler godoc
//
//	@Summary		Get Orders
//	@Description	get all orders of an user
//	@Tags			Get Orders
//	@Accept			json
//	@Produce		json
//	@Param			user	query		string	true	"user"
//	@Success		200		{object}	int
//	@Failure		404		{object}	string	"page not found"
//	@Router			/getorders/ [get]
func (hc *HandlerCore) GetOrdersHandler(c *gin.Context) {
	userAddr := c.Query("user")

	// connect to an eth node with ep
	logger.Info("connecting chain")
	backend, chainID := eth.ConnETH(Chain_Endpoint)
	logger.Info("chain id:", chainID)

	// get contract instance
	marketIns, err := market.NewMarket(common.HexToAddress(comm.Contracts.Market), backend)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("new market instance failed: %s", err.Error())})
	}

	// get orders of an user
	orders, err := marketIns.GetOrders(&bind.CallOpts{}, common.HexToAddress(userAddr))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("get orders failed: %s", err.Error())})
	}

	// response order list
	c.JSON(http.StatusOK, orders)

}

// get global info
// CpsHandler godoc
//
//	@Summary		Get Global Info
//	@Description	get global info
//	@Tags			Get GInfo
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	int
//	@Failure		404	{object}	string	"page not found"
//	@Router			/ginfo/ [get]
func (hc *HandlerCore) GInfoHandler(c *gin.Context) {
	// connect to an eth node with ep
	logger.Info("connecting chain")
	backend, chainID := eth.ConnETH(Chain_Endpoint)
	logger.Info("chain id:", chainID)

	logger.Info("get registry instance")
	// get contract instance
	regIns, err := registry.NewRegistry(common.HexToAddress(comm.Contracts.Registry), backend)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("new market instance failed: %s", err.Error())})
		return
	}

	logger.Info("getting cp number")

	//
	gInfo, err := regIns.GInfo(&bind.CallOpts{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("call contract failed: %s", err.Error())})
		return
	}
	//logger.Info("cp keys:", keys)

	// response
	c.JSON(http.StatusOK, gInfo)
}

// qeury credit for a role with address
//
//	@Summary		QueryCredit
//	@Description	Query credit of someone
//	@Tags			QueryCredit
//	@Accept			json
//	@Produce		json
//	@Param			address	query		string	true	"address of this caller"
//	@Success		200		{object}	string	"query OK"
//	@Failure		400		{object}	string	"bad request"
//	@Router			/querycredit [get]
func (hc *HandlerCore) QueryCreditHandler(c *gin.Context) {
	userAddr := c.Query("address")
	creditAddr := comm.Contracts.Credit

	// connect to an eth node with ep
	backend, chainID := eth.ConnETH(Chain_Endpoint)
	fmt.Println("chain id:", chainID)

	// get credit instance
	creditIns, err := credit.NewCredit(common.HexToAddress(creditAddr), backend)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// query balance
	bal, err := creditIns.BalanceOf(&bind.CallOpts{}, common.HexToAddress(userAddr))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"credit balance": bal})

}

// get current version of contracts

// handler of version
//
//	@Summary		version
//	@Description	get version
//	@Tags			Version
//	@Accept			multipart/form-data
//	@Produce		json
//	@Success		200	{object}	string	"version OK"
//	@Failure		400	{object}	string	"bad request"
//	@Router			/version [get]
func (hc *HandlerCore) CurrentVerHandler(c *gin.Context) {
	verAddr := comm.Contracts.Version

	// connect to an eth node with ep
	backend, chainID := eth.ConnETH(Chain_Endpoint)
	fmt.Println("chain id:", chainID)

	// version instance
	verIns, err := version.NewVersion(common.HexToAddress(verAddr), backend)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// call version
	ver, err := verIns.CurrentVer(&bind.CallOpts{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"current contracts version": ver})

}

// for cross domain
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
