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
	comm "github.com/rockiecn/platform/common"
	"github.com/rockiecn/platform/lib/kv"
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
//
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
//	@Param			numDisk	formData	string	true	"num"
//	@Param			priDisk	formData	string	false	"price"
//	@Param			numMem		formData	string	true	"num"
//	@Param			priMem		formData	string	false	"price"
//	@Success		200			{object}	string	"regist OK"
//	@Failure		400			{object}	string	"bad request"
//	@Router			/registcp [post]
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
	client, err := ethclient.Dial(eth.Endpoint)
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
	eth.CheckTx(eth.Endpoint, signedTx.Hash(), "")

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
	client, err := ethclient.Dial(eth.Endpoint)
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
	eth.CheckTx(eth.Endpoint, signedTx.Hash(), "")

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
//	@Success		200	{object}	[]CPInfo
//	@Failure		404	{object}	string	"page not found"
//
//	@Router			/listcp/ [get]
func (hc *HandlerCore) ListCPHandler(c *gin.Context) {

	// connect to an eth node with ep
	logger.Info("connecting chain")
	backend, chainID := eth.ConnETH(eth.Endpoint)
	logger.Info("chain id:", chainID)

	logger.Info("registry address: ", comm.Contracts.Registry)
	// get contract instance
	registryIns, err := registry.NewRegistry(common.HexToAddress(comm.Contracts.Registry), backend)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("new registry instance failed: %s", err.Error()).Error()})
		return
	}

	// get cp list
	list, err := registryIns.GetList(&bind.CallOpts{})
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
//	@Tags			get cp
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
	backend, chainID := eth.ConnETH(eth.Endpoint)
	logger.Info("chain id:", chainID)

	// get contract instance
	contractIns, err := registry.NewRegistry(common.HexToAddress(comm.Contracts.Registry), backend)

	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("new token instance failed: %v", err))
		return
	}

	// get balance of addr2
	regInfo, err := contractIns.Get(&bind.CallOpts{}, common.HexToAddress(cpaddr))
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
	client, err := ethclient.Dial(eth.Endpoint)
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
	eth.CheckTx(eth.Endpoint, signedTx.Hash(), "")

	// response
	c.JSON(http.StatusOK, gin.H{
		"msg": "user approve ok",
	})
}

// check credit allowance
func (hc *HandlerCore) AllowanceHandler(c *gin.Context) {
	// param in form
	owner := c.Query("owner")
	spender := c.Query("spender")

	// connect to an eth node with ep
	logger.Info("connecting chain")
	backend, chainID := eth.ConnETH(eth.Endpoint)
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
//
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
//	@Param			numDisk	formData	string	true	"num"
//	@Param			priDisk	formData	string	false	"price"
//	@Param			numMem		formData	string	true	"num"
//	@Param			priMem		formData	string	false	"price"
//	@Param			duration	formData	string	true	"duration"
//	@Success		200			{object}	string	"regist OK"
//	@Failure		400			{object}	string	"bad request"
//	@Router			/createorder [post]
func (hc *HandlerCore) CreateOrderHandler(c *gin.Context) {
	// tx data in form
	txData := c.PostForm("tx")

	// transfer to types.Transaction
	signedTx := new(types.Transaction)
	signedTx.UnmarshalJSON([]byte(txData))
	log.Println("signed tx: ", signedTx)

	// connect to an eth client
	log.Println("connecting client")
	client, err := ethclient.Dial(eth.Endpoint)
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
	eth.CheckTx(eth.Endpoint, signedTx.Hash(), "")

	// response
	c.JSON(http.StatusOK, gin.H{
		"msg": "create order ok",
	})
}

// handler for getOrder
func (hc *HandlerCore) GetOrderHandler(c *gin.Context) {

	// user and cp
	userAddr := c.Query("user")
	cpAddr := c.Query("cp")

	// connect to an eth node with ep
	logger.Info("connecting chain")
	backend, chainID := eth.ConnETH(eth.Endpoint)
	logger.Info("chain id:", chainID)

	// get contract instance
	marketIns, err := market.NewMarket(common.HexToAddress(comm.Contracts.Market), backend)

	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("new market instance failed: %v", err))
		return
	}

	// get order with user and cp
	orderInfo, err := marketIns.GetOrder(&bind.CallOpts{From: eth.Addr1}, common.HexToAddress(userAddr), common.HexToAddress(cpAddr))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logger.Info("order info:", orderInfo)

	// response node list
	c.JSON(http.StatusOK, orderInfo)
}

// handler for get order list for user or cp
//
//	@Summary		List order
//	@Description	list an order
//	@Tags			ListOrder
//	@Accept			json
//	@Produce		json
//	@Param			role	query		string	true	"user or provider"
//	@Param			address	query		string	true	"address"
//
//	@Success		200		{object}	string	"list OK"
//	@Failure		400		{object}	string	"bad request"
//	@Router			/listorder [get]
// func (hc *HandlerCore) ListOrderHandler(c *gin.Context) {

// 	// connect to an eth node with ep
// 	logger.Info("connecting chain")
// 	backend, chainID := eth.ConnETH(eth.Endpoint)
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

// qeury credit for a role with address
//
//	@Summary		QueryCredit
//	@Description	Query credit of someone
//	@Tags			QueryCredit
//	@Accept			json
//	@Produce		json
//	@Param			role	query		string	true	"role of this caller"
//	@Param			address	query		string	true	"address of this caller"
//	@Success		200		{object}	string	"query OK"
//	@Failure		400		{object}	string	"bad request"
//	@Router			/querycredit [get]
func (hc *HandlerCore) QueryCreditHandler(c *gin.Context) {
	userAddr := c.Query("address")
	creditAddr := comm.Contracts.Credit

	// connect to an eth node with ep
	backend, chainID := eth.ConnETH(eth.Endpoint)
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
