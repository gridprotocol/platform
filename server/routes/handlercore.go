package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
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

	// get contract instance
	registryIns, err := registry.NewRegistry(common.HexToAddress(comm.Contracts.Registry), backend)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("new registry instance failed: %s", err.Error())})
	}

	// get key list
	keys, err := registryIns.GetKeys(&bind.CallOpts{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("get cp keys failed: %s", err.Error())})
	}
	//logger.Info("cp keys:", keys)

	// response key list
	c.JSON(http.StatusOK, keys)
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
	}

	// get balance of addr2
	regInfo, err := contractIns.Get(&bind.CallOpts{}, common.HexToAddress(cpaddr))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	logger.Info("registry info:", regInfo)

	// response node list
	c.JSON(http.StatusOK, regInfo)
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

	// mint credit to user and user approving

	// connect to an eth node with ep
	backend, chainID := eth.ConnETH(eth.Endpoint)
	fmt.Println("chain id:", chainID)

	// auth for admin
	authAdmin, err := eth.MakeAuth(chainID, eth.SK0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// credit instance
	creditIns, err := credit.NewCredit(common.HexToAddress(comm.Contracts.Credit), backend)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// total value of test order
	totalValue, ok := new(big.Int).SetString("2831300", 10)
	if !ok {
		log.Fatal("big set string failed")
	}

	// mint some credit for approve
	fmt.Println("admin mint some credit to user for create order")
	tx, err := creditIns.Mint(authAdmin, eth.Addr1, totalValue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("waiting for tx to be ok")
	err = eth.CheckTx(eth.Endpoint, tx.Hash(), "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// approve must be done by the user before create an order
	/*
		// auth for user
		authUser, err := eth.MakeAuth(chainID, eth.SK1)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fmt.Println("user approving credit to market.., approve value: ", totalValue)
		tx, err = creditIns.Approve(authUser, common.HexToAddress(comm.Contracts.Market), totalValue)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// wait for tx to be ok
		fmt.Println("waiting tx")
		err = eth.CheckTx(eth.Endpoint, tx.Hash(), "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	*/

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
	cpaddr := c.Query("cp")

	// connect to an eth node with ep
	logger.Info("connecting chain")
	backend, chainID := eth.ConnETH(eth.Endpoint)
	logger.Info("chain id:", chainID)

	// get contract instance
	marketIns, err := market.NewMarket(common.HexToAddress(comm.Contracts.Market), backend)

	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("new market instance failed: %v", err))
	}

	// get order with user and cp
	orderInfo, err := marketIns.GetOrder(&bind.CallOpts{From: common.HexToAddress(userAddr)}, common.HexToAddress(cpaddr))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
func (hc *HandlerCore) ListOrderHandler(c *gin.Context) {

	// get role
	role := c.Query("role")
	// user address from param
	addr := c.Query("address")

	var orderList []OrderInfo
	var err error

	// order list for response
	switch role {
	case "user":
		orderList, err = hc.getUserOrders(addr)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
	case "cp":
		orderList, err = hc.getCpOrders(addr)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
	default:
		c.JSON(http.StatusOK, gin.H{"response": "error type in request"})
	}

	// response order list
	c.JSON(http.StatusOK, orderList)
}

// user record credit with txHash
// value - uint: eth
// credit = eth * 1000000
//
//	@Summary		Pay for credit
//	@Description	Pay to credit with a transfer's key
//	@Tags			Pay
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			transkey	formData	string	true	"transfer key"
//	@Success		200			{object}	string	"pay OK"
//	@Failure		400			{object}	string	"bad request"
//	@Router			/pay [post]
func (hc *HandlerCore) PayHandler(c *gin.Context) {
	// get key of a transfer
	transkey := c.PostForm("transkey")
	transfer := TransferInfo{}
	// get transfer info
	data, err := hc.LocalDB.Get([]byte(transkey))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	// unmarshal transfer info
	err = json.Unmarshal(data, &transfer)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	// read field
	from := transfer.From
	value := transfer.Value
	txHash := transfer.TxHash
	confirmed := transfer.TxConfirmed
	saved := transfer.CreditSaved

	// tx not confirmed
	if !confirmed {
		c.JSON(http.StatusOK, gin.H{"error": "tx of this transfer has not been confirmed on chain yet"})
		return
	}

	// transfer already used
	if saved {
		c.JSON(http.StatusOK, gin.H{"error": "this transfer has been used for credit"})
		return
	}

	value64, err := utils.StringToInt64(value)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	// calc credit: eth * 10^6
	credit := value64 * 1000000

	// get credit
	oldCredit, err := hc.queryCredit(from)
	if err != nil {
		if err.Error() == "Key not found" {
			oldCredit = "0"
		} else {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
	}

	logger.Debug("old credit:", oldCredit)

	// accumulate credit
	old64, err := utils.StringToInt64(oldCredit)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	new64 := old64 + credit
	new := utils.Int64ToString(new64)

	logger.Debug("new credit:", new)

	// for atomic
	keys := [][]byte{}
	values := [][]byte{}

	// update credit for this account
	creKey := CreditKey(from)
	keys = append(keys, creKey)
	values = append(values, []byte(new))

	// get payinfo id for this account
	oldID, err := hc.getPayInfoID(from)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	logger.Debug("old credit id:", oldID)

	// update payinfo id
	oldID64, err := utils.StringToInt64(oldID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	newID := utils.Int64ToString(oldID64 + 1)
	logger.Debug("new payinfo id:", newID)
	// update payinfo id for this account
	idKey := PayInfoIDKey(from)
	keys = append(keys, idKey)
	values = append(values, []byte(newID))

	// make payinfo's key
	piKey := PayInfoKey(from, oldID)
	logger.Debugf("payinfo key:%s", piKey)

	// record pay info into db
	payInfo := PayInfo{
		PIKey:  string(piKey),
		TIKey:  transkey, // which transfer is used for this credit
		Owner:  from,
		Credit: credit,
		TxHash: txHash,
	}
	// marshal pi to bytes
	data, err = json.Marshal(payInfo)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	// record payinfo data
	keys = append(keys, piKey)
	values = append(values, data)

	// modify transferinfo's state(tx confirmed, credit saved)
	transfer.TxConfirmed = true
	transfer.CreditSaved = true

	// marshal to bytes
	data, err = json.Marshal(transfer)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	keys = append(keys, []byte(transkey))
	values = append(values, data)

	// multiput
	hc.LocalDB.MultiPut(keys, values)

	// response
	c.JSON(http.StatusOK, gin.H{"response": "pay ok"})
}

// query pay infos
//
//	@Summary		ListPay
//	@Description	ListPay
//	@Tags			ListPay
//	@Accept			json
//	@Produce		json
//	@Param			addr	query		string	true	"address of an user"
//	@Success		200		{object}	string	"list pay OK"
//	@Failure		400		{object}	string	"bad request"
//	@Router			/listpay [get]
func (hc *HandlerCore) ListPayHandler(c *gin.Context) {
	addr := c.Query("addr")

	piList, err := hc.getPayInfoList(addr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	// response
	c.JSON(http.StatusOK, piList)
}

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
		fmt.Println("new credit instance failed:", err)
	}

	// query balance
	bal, err := creditIns.BalanceOf(&bind.CallOpts{}, common.HexToAddress(userAddr))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"credit balance": bal})

}

// transfer, write transfer records into db
// transfer info key: trans_user_id
// id key: trans_*
//
//	@Summary		Transfer token
//	@Description	user transfer token to platform
//	@Tags			Transfer
//	@Accept			json
//	@Produce		json
//	@Param			txHash	query		string	true	"tx hash"
//	@Param			from	query		string	true	"from addr"
//	@Param			to		query		string	true	"to addr"
//	@Param			value	query		string	true	"transfer value"
//	@Success		200		{object}	string	"transfer OK"
//	@Failure		400		{object}	string	"bad request"
//	@Router			/transfer [post]
func (hc *HandlerCore) TransferHandler(c *gin.Context) {
	txHash := c.Query("txHash")
	from := c.Query("from")
	to := c.Query("to")
	value := c.Query("value")

	confirmed := false
	creditSaved := false

	// for atomic
	keys := [][]byte{}
	values := [][]byte{}

	// for id update
	id, err := hc.getTransferID(from)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	id64, err := utils.StringToInt64(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	id64++
	newID := utils.Int64ToString(id64)
	// transfer id key
	idKey := TransferIDKey(from)
	// update transfer id
	keys = append(keys, []byte(idKey))
	values = append(values, []byte(newID))

	// key for transfer info
	tiKey := TransferInfoKey(from, id)
	// make ti
	ti := TransferInfo{
		TIKey:       string(tiKey),
		TxHash:      txHash,
		From:        from,
		To:          to,
		Value:       value,
		TxConfirmed: confirmed,
		CreditSaved: creditSaved,
	}
	data, err := json.Marshal(ti)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	// record ti
	keys = append(keys, tiKey)
	values = append(values, data)

	// multiput
	hc.LocalDB.MultiPut(keys, values)

	c.JSON(http.StatusOK, gin.H{"response": "transfer ok"})
}

// list all transfer info about an user
//
//	@Summary		List all transfers
//	@Description	List all transfers of an address
//	@Tags			List transfers
//	@Accept			json
//	@Produce		json
//	@Param			address	query		string	true	"address to show list"
//	@Success		200		{object}	string	"list transfer OK"
//	@Failure		400		{object}	string	"bad request"
//	@Router			/listtransfer [get]
func (hc *HandlerCore) ListTransferHandler(c *gin.Context) {

	// user address from param
	addr := c.Query("address")

	// transfer list for response
	transList, err := hc.getUserTransfers(addr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transList)
}

// refresh all transfers of an user, check transfers' confirmed state
//
//	@Summary		RefreshTransfer status of transfer
//	@Description	Refresh status of transfer of an address
//	@Tags			Refresh Transfer
//	@Accept			json
//	@Produce		json
//	@Param			address	query		string	true	"address to refresh"
//	@Success		200		{object}	string	"refresh OK"
//	@Failure		400		{object}	string	"bad request"
//	@Router			/refreshtransfer [post]
func (hc *HandlerCore) RefreshTransferHandler(c *gin.Context) {
	userAddr := c.Query("address")

	transfers, err := hc.getUserTransfers(userAddr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	// for atomic
	keys := [][]byte{}
	values := [][]byte{}

	// check all transfers for confirm
	for _, transfer := range transfers {
		// if this transfer already confirmed, skip
		if transfer.TxConfirmed {
			continue
		}

		// check for now
		confirmed, err := checkTxConfirmed(transfer.TxHash)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		if confirmed {
			k, v, err := hc.setTransferConfirmed([]byte(transfer.TIKey), true)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
				return
			}
			// k,v
			keys = append(keys, k)
			values = append(values, v)
		}
	}

	// multi put all k,v
	hc.LocalDB.MultiPut(keys, values)

	c.JSON(http.StatusOK, gin.H{"response": "refresh transfer ok"})
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
