package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgraph-io/badger/v2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/grid/contracts/eth"
	"github.com/grid/contracts/go/registry"
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

var ctr = eth.Address{}

// load all contract addresses from json
func init() {
	ctr = eth.LoadLocal()
	logger.Info("contract addresses:", ctr)
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

	// provider name
	//name := c.PostForm("name")

	// provider wallet address
	cpAddr := c.PostForm("cp")
	_ = cpAddr

	// connection info of provider
	ip := c.PostForm("ip")
	domain := c.PostForm("domain")
	port := c.PostForm("port")

	// resource info
	numCPU := c.PostForm("nCPU")
	priCPU := c.PostForm("pCPU")

	numGPU := c.PostForm("nGPU")
	priGPU := c.PostForm("pGPU")

	numDisk := c.PostForm("nDisk")
	priDisk := c.PostForm("pDisk")

	numMem := c.PostForm("nMem")
	priMem := c.PostForm("pMem")

	// check input
	if !isNumber(priCPU) {
		c.JSON(http.StatusBadRequest, gin.H{"response": "cpu price must be number"})
		return
	}
	if !isNumber(priGPU) {
		c.JSON(http.StatusBadRequest, gin.H{"response": "gpu price must be number"})
		return
	}
	if !isNumber(priMem) {
		c.JSON(http.StatusBadRequest, gin.H{"response": "mem price must be number"})
		return
	}
	if !isNumber(priDisk) {
		c.JSON(http.StatusBadRequest, gin.H{"response": "store price must be number"})
		return
	}
	if !isNumber(numDisk) {
		c.JSON(http.StatusBadRequest, gin.H{"response": "store space must be number"})
		return
	}
	if !isNumber(numMem) {
		c.JSON(http.StatusBadRequest, gin.H{"response": "memory space must be number"})
		return
	}
	if !isNumber(numCPU) {
		c.JSON(http.StatusBadRequest, gin.H{"response": "cpu number must be number"})
		return
	}
	if !isNumber(numGPU) {
		c.JSON(http.StatusBadRequest, gin.H{"response": "gpu number must be number"})
		return
	}

	// cal set of registry contract to register a cp info

	// connect to an eth node with ep
	backend, chainID := eth.ConnETH(eth.Endpoint)
	logger.Info("chain id:", chainID)

	// get registry instance
	regIns, err := registry.NewRegistry(common.HexToAddress(ctr.Registry), backend)
	if err != nil {
		logger.Error("new registry instance failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("making auth")
	// admin register for cp
	// todo: cp need pay credit for tx fee
	txAuth, err := eth.MakeAuth(chainID, eth.SK1)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logger.Info("make auth ok")

	// provider as register
	//txAuth.From = common.HexToAddress(cpAddr)

	// type transfer
	cpu, err := utils.StringToUint64(numCPU)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	gpu, err := utils.StringToUint64(numGPU)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	mem, err := utils.StringToUint64(numMem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	disk, err := utils.StringToUint64(numDisk)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	logger.Info("call registry.set")

	// call registry's Set method with all params
	tx, err := regIns.Set(txAuth, ip, domain, port, cpu, gpu, mem, disk)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// wait tx ok
	logger.Info("waiting for set to be ok")
	eth.CheckTx(eth.Endpoint, tx.Hash(), "")

	// get cp's reg info
	regInfo, err := regIns.Get(&bind.CallOpts{}, eth.Addr1)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("regInfo:", regInfo)

	// response to client
	c.JSON(http.StatusOK, gin.H{"response": "regist OK"})

	/*
		//
		info := CPInfo{
			Name:     name,
			Address:  address,
			EndPoint: endpoint,
			NumCPU:   numCPU,
			PriCPU:   priCPU,
			NumGPU:   numGPU,
			PriGPU:   priGPU,
			numDisk: numDisk,
			priDisk: priDisk,
			NumMem:   numMem,
			PriMem:   priMem,
		}

		// marshal into bytes
		data, err := json.Marshal(info)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		KEY := CPInfoKey(address)

		// check if cp exists
		b, err := hc.LocalDB.Has(KEY)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		if b {
			c.JSON(http.StatusOK, gin.H{"response": "cp already exist"})
			return
		}

		// wallet address as key, info as valude
		err = hc.LocalDB.Put(KEY, data)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"response": "regist OK"})
	*/
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
	// all cp info to response
	cps := make([]CPInfo, 0, 100)

	// get all cp info with prefix
	PREFIX := []byte("CP_INFO_") // 设置通配符前缀
	err := hc.LocalDB.DB.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(PREFIX); it.ValidForPrefix(PREFIX); it.Next() {
			err := hc.appendResult(&cps, it.Item())
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	// response node list
	c.JSON(http.StatusOK, cps)
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
	addr := c.Query("address")

	cp := CPInfo{}

	// get cpinfo key
	cpkey := CPInfoKey(addr)

	// read db for cp info
	v, err := hc.LocalDB.Get(cpkey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// unmarshal
	err = json.Unmarshal(v, &cp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Debug("cpinfo:", cp)

	// response node list
	c.JSON(http.StatusOK, cp)
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

	// user address
	userAddr := c.PostForm("userAddress")

	// provider address
	cpAddr := c.PostForm("cpAddress")

	numCPU := c.PostForm("numCPU")
	priCPU := c.PostForm("priCPU")

	numGPU := c.PostForm("numGPU")
	priGPU := c.PostForm("priGPU")

	numDisk := c.PostForm("numDisk")
	priDisk := c.PostForm("priDisk")

	numMem := c.PostForm("numMem")
	priMem := c.PostForm("priMem")

	// duration in month
	dur := c.PostForm("duration")

	// check input
	if !isNumber(priCPU) {
		c.JSON(http.StatusBadRequest, gin.H{"response": "price must be number"})
		return
	}
	if !isNumber(priGPU) {
		c.JSON(http.StatusBadRequest, gin.H{"response": "price must be number"})
		return
	}
	if !isNumber(priMem) {
		c.JSON(http.StatusBadRequest, gin.H{"response": "price must be number"})
		return
	}
	if !isNumber(priDisk) {
		c.JSON(http.StatusBadRequest, gin.H{"response": "price must be number"})
		return
	}
	if !isNumber(numDisk) {
		c.JSON(http.StatusBadRequest, gin.H{"response": "store space must be number"})
		return
	}
	if !isNumber(numMem) {
		c.JSON(http.StatusBadRequest, gin.H{"response": "memory space must be number"})
		return
	}

	// compute expire with duration and current time
	expire, err := utils.DurToTS(dur)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	// read cp name from db with cp address

	// get cp info key
	cpkey := CPInfoKey(cpAddr)
	// check cp
	b, err := hc.LocalDB.Has([]byte(cpkey))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	if !b {
		c.JSON(http.StatusOK, gin.H{"response": "cp not found"})
		return
	}
	// read cp info
	data, err := hc.LocalDB.Get([]byte(cpkey))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	// unmarshal data to CPInfo
	cp := CPInfo{}
	err = json.Unmarshal(data, &cp)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	// get cp name and endpoint
	cpName := cp.Name
	endPoint := cp.EndPoint

	// get current order id for each user, used in new order
	var orderID string
	orderID, err = hc.getOrderID(userAddr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	logger.Debugf("old order id:%s", orderID)

	// order key
	orderKey := OrderKey(userAddr, orderID)
	logger.Debugf("order key:%s", orderKey)

	// construct new order info
	order := OrderInfo{
		OrderKey: string(orderKey),
		UserAddr: userAddr,
		CPAddr:   cpAddr,
		CPName:   cpName,
		EndPoint: endPoint,
		NumCPU:   numCPU,
		PriCPU:   priCPU,
		NumGPU:   numGPU,
		PriGPU:   priGPU,
		NumDisk:  numDisk,
		PriDisk:  priDisk,

		NumMem:  numMem,
		PriMem:  priMem,
		Dur:     dur,
		Expire:  expire,
		Settled: false,
	}

	logger.Debug("GPU price:", order.PriGPU)

	// calc credit cost of order
	cost64, err := CalcCost(&order)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	logger.Debug("credit cost:", cost64)

	// set cost in order
	order.Cost = cost64

	// get credit
	credit, err := hc.queryCredit(userAddr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	logger.Debug("credit left:", credit)

	// check credit
	credit64, err := utils.StringToInt64(credit)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	if credit64 < cost64 {
		c.JSON(
			http.StatusOK,
			gin.H{"response": "credit is not enough to pay this order,create order failed"},
		)
		return
	}

	// for atomic operations on db
	keys := [][]byte{}
	values := [][]byte{}

	// update user's credit
	credit64 = credit64 - cost64
	newCredit := utils.Int64ToString(credit64)
	logger.Debug("new credit after createorder:", newCredit)

	// update user's credit in db
	creKey := CreditKey(userAddr)
	keys = append(keys, []byte(creKey))
	values = append(values, []byte(newCredit))

	// db operation

	// mashal order into bytes
	data, err = json.Marshal(order)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	// put order info into db
	keys = append(keys, orderKey)
	values = append(values, data)

	// increase order id by 1
	orderID64, err := utils.StringToInt64(orderID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	orderID64 += 1
	orderID = utils.Int64ToString(orderID64)
	logger.Debugf("new order id:%s", orderID)
	// update order id
	idKey := OrderIDKey(userAddr)
	keys = append(keys, idKey)
	values = append(values, []byte(orderID))

	// append the order key for cp into db
	k, v, err := hc.appendOrder(cpAddr, string(orderKey))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	keys = append(keys, k)
	values = append(values, v)

	// for atomic
	err = hc.LocalDB.MultiPut(keys, values)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	// response
	c.JSON(http.StatusOK, gin.H{
		"cost": cost64, // credit = eth*1000000
	})
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
	role := c.Query("role")
	address := c.Query("address")

	var credit string

	switch role {
	case "user":
		// get old credit from db, if key not found, init with 0
		credit, err := hc.queryCredit(address)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		logger.Debug("credit:", credit)

		c.JSON(http.StatusOK, gin.H{
			"credit": credit,
		})

	case "cp":
		// settle all orders of this cp, set order state

		// get cp's order list
		orderList, err := hc.getCpOrders(address)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		// deal with each order
		for _, o := range orderList {
			// for multiput
			keys := [][]byte{}
			values := [][]byte{}

			// get current time stamp
			now := time.Now().Unix()
			logger.Debug("current timestamp:", now)
			expire := o.Expire
			expire64, err := utils.StringToInt64(expire)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
				return
			}
			logger.Debug("expire timestamp:", expire64)

			// todo: if order not expired, skip it
			// if expire64 < now {
			// 	continue
			// }

			// if not settled
			if !o.Settled {
				// add order's cost into cp's credit
				k, v, err := hc.addCredit(o.CPAddr, o.Cost)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{"error": err.Error()})
					return
				}
				keys = append(keys, k)
				values = append(values, v)

				// set order's settled state to true
				k, v, err = hc.setOrderSettled([]byte(o.OrderKey), true)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{"error": err.Error()})
					return
				}
				keys = append(keys, k)
				values = append(values, v)

				// multiput
				hc.LocalDB.MultiPut(keys, values)
			}
		}

		// get credit from db, if key not found, response 0
		credit, err = hc.queryCredit(address)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		logger.Debug("credit:", credit)

		// response credit
		c.JSON(http.StatusOK, gin.H{
			"credit": credit,
		})
	default:
		c.JSON(http.StatusOK, gin.H{"response": "error role in request"})
	}
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
