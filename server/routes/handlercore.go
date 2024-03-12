package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgraph-io/badger/v2"
	"github.com/gin-gonic/gin"
	"github.com/rockiecn/platform/lib/kv"
	"github.com/rockiecn/platform/lib/utils"
)

type HandlerCore struct {
	LocalDB *kv.Database
}

// pay info when recharge credit, need to be stored in db
type PayInfo struct {
	PIKey  string `json:"PayInfoKey"`
	Payer  string `json:"Payer"`
	Credit uint64 `json:"Credit"`
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

// handler of welcom
func (hc *HandlerCore) RootHandler(c *gin.Context) {
	c.String(http.StatusOK, "Welcome Server")
}

// handler of cp login
func (hc *HandlerCore) RegistCPHandler(c *gin.Context) {

	// provider name
	name := c.PostForm("name")
	// provider wallet address
	address := c.PostForm("address")

	endpoint := c.PostForm("endpoint")

	numCPU := c.PostForm("numCPU")
	priCPU := c.PostForm("priCPU")

	numGPU := c.PostForm("numGPU")
	priGPU := c.PostForm("priGPU")

	numStore := c.PostForm("numStore")
	priStore := c.PostForm("priStore")

	numMem := c.PostForm("numMem")
	priMem := c.PostForm("priMem")

	// check input
	if !isNumber(priCPU) {
		c.JSON(http.StatusBadRequest, "price must be number")
		return
	}
	if !isNumber(priGPU) {
		c.JSON(http.StatusBadRequest, "price must be number")
		return
	}
	if !isNumber(priMem) {
		c.JSON(http.StatusBadRequest, "price must be number")
		return
	}
	if !isNumber(priStore) {
		c.JSON(http.StatusBadRequest, "price must be number")
		return
	}
	if !isNumber(numStore) {
		c.JSON(http.StatusBadRequest, "store space must be number")
		return
	}
	if !isNumber(numMem) {
		c.JSON(http.StatusBadRequest, "memory space must be number")
		return
	}

	//
	info := CPInfo{
		Name:     name,
		Address:  address,
		EndPoint: endpoint,
		NumCPU:   numCPU,
		PriCPU:   priCPU,
		NumGPU:   numGPU,
		PriGPU:   priGPU,
		NumStore: numStore,
		PriStore: priStore,
		NumMem:   numMem,
		PriMem:   priMem,
	}

	// marshal into bytes
	data, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}

	KEY := CPInfoKey(address)

	// check if cp exists
	b, err := hc.LocalDB.Has(KEY)
	if err != nil {
		panic(err)
	}
	if b {
		c.JSON(http.StatusOK, "cp already exist")
		return
	}

	// wallet address as key, info as valude
	err = hc.LocalDB.Put(KEY, data)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, "regist OK")
}

// handler for list cp nodes
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
		panic(err)
	}

	// response node list
	c.JSON(http.StatusOK, cps)
}

// handler of create order
func (hc *HandlerCore) CreateOrderHandler(c *gin.Context) {

	// user address
	userAddr := c.PostForm("userAddress")

	// provider address
	cpAddr := c.PostForm("cpAddress")

	numCPU := c.PostForm("numCPU")
	priCPU := c.PostForm("priCPU")

	numGPU := c.PostForm("numGPU")
	priGPU := c.PostForm("priGPU")

	numStore := c.PostForm("numStore")
	priStore := c.PostForm("priStore")

	numMem := c.PostForm("numMem")
	priMem := c.PostForm("priMem")

	// duration in month
	dur := c.PostForm("duration")

	// check input
	if !isNumber(priCPU) {
		c.JSON(http.StatusBadRequest, "price must be number")
		return
	}
	if !isNumber(priGPU) {
		c.JSON(http.StatusBadRequest, "price must be number")
		return
	}
	if !isNumber(priMem) {
		c.JSON(http.StatusBadRequest, "price must be number")
		return
	}
	if !isNumber(priStore) {
		c.JSON(http.StatusBadRequest, "price must be number")
		return
	}
	if !isNumber(numStore) {
		c.JSON(http.StatusBadRequest, "store space must be number")
		return
	}
	if !isNumber(numMem) {
		c.JSON(http.StatusBadRequest, "memory space must be number")
		return
	}

	// compute expire with duration and current time
	expire, err := utils.DurToTS(dur)
	if err != nil {
		panic(err)
	}

	// read cp name from db with cp address

	// get cp info key
	cpkey := CPInfoKey(cpAddr)
	// check cp
	b, err := hc.LocalDB.Has([]byte(cpkey))
	if err != nil {
		panic(err)
	}
	if !b {
		c.JSON(http.StatusOK, "cp not found")
		return
	}
	// read cp info
	data, err := hc.LocalDB.Get([]byte(cpkey))
	if err != nil {
		panic(err)
	}
	// unmarshal data to CPInfo
	cp := CPInfo{}
	err = json.Unmarshal(data, &cp)
	if err != nil {
		panic(err)
	}
	// get cp name and endpoint
	cpName := cp.Name
	endPoint := cp.EndPoint

	// get current order id for each user, used in new order
	var orderID string
	orderID, err = hc.getOrderID(userAddr)
	if err != nil {
		panic(err)
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
		NumStore: numStore,
		PriStore: priStore,
		NumMem:   numMem,
		PriMem:   priMem,
		Dur:      dur,
		Expire:   expire,
		Settled:  false,
	}

	// calc credit cost of order
	cost64, err := CalcCost(&order)
	if err != nil {
		panic(err)
	}
	logger.Debug("credit cost:", cost64)

	// set cost in order
	order.Cost = cost64

	// get credit
	credit, err := hc.getCredit(userAddr)
	if err != nil {
		panic(err)
	}

	logger.Debug("credit left:", credit)

	// check credit
	credit64, err := utils.StringToUint64(credit)
	if err != nil {
		panic(err)
	}
	if credit64 < cost64 {
		c.JSON(http.StatusOK, "credit is not enough to pay this order,create order failed")
		return
	}

	// for atomic operations on db
	keys := [][]byte{}
	values := [][]byte{}

	// update user's credit
	credit64 = credit64 - cost64
	newCredit := utils.Uint64ToString(credit64)
	logger.Debug("new credit after createorder:", newCredit)

	// update user's credit in db
	creKey := CreditKey(userAddr)
	keys = append(keys, []byte(creKey))
	values = append(values, []byte(newCredit))

	// db operation

	// mashal order into bytes
	data, err = json.Marshal(order)
	if err != nil {
		panic(err)
	}
	// put order info into db
	keys = append(keys, orderKey)
	values = append(values, data)

	// increase order id by 1
	orderID64, err := utils.StringToUint64(orderID)
	if err != nil {
		panic(err)
	}
	orderID64 += 1
	orderID = utils.Uint64ToString(orderID64)
	logger.Debugf("new order id:%s", orderID)
	// update order id
	idKey := OrderIDKey(userAddr)
	keys = append(keys, idKey)
	values = append(values, []byte(orderID))

	// append the order key for cp into db
	k, v, err := hc.appendOrder(cpAddr, string(orderKey))
	if err != nil {
		panic(err)
	}
	keys = append(keys, k)
	values = append(values, v)

	// for atomic
	err = hc.LocalDB.MultiPut(keys, values)
	if err != nil {
		panic(err)
	}

	// response
	c.JSON(http.StatusOK, gin.H{
		//"to":   "0x1234",
		"cost": cost64, // credit = eth*1000000
	})
}

// handler for get order list for user or cp
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
			panic(err)
		}
	case "cp":
		orderList, err = hc.getCpOrders(addr)
		if err != nil {
			panic(err)
		}
	default:
		c.JSON(http.StatusOK, "error type in request")
	}

	// response order list
	c.JSON(http.StatusOK, orderList)
}

// user record credit with txHash
// value uint: eth
// crecit = eth * 1000000
func (hc *HandlerCore) PayHandler(c *gin.Context) {
	from := c.PostForm("from")
	//to := c.PostForm("to")
	value := c.PostForm("value")
	txHash := c.PostForm("txHash")

	//todo: verify txHash
	txConfirmed := false
	_ = txConfirmed

	value64, err := utils.StringToUint64(value)
	if err != nil {
		panic(err)
	}
	// calc credit: eth * 10^6
	credit := value64 * 1000000

	// get credit
	oldCredit, err := hc.getCredit(from)
	if err != nil {
		if err.Error() == "Key not found" {
			oldCredit = "0"
		} else {
			panic(err)
		}
	}

	logger.Debug("old credit:", oldCredit)

	// accumulate credit
	old64, err := utils.StringToUint64(oldCredit)
	if err != nil {
		panic(err)
	}
	new64 := old64 + credit
	new := utils.Uint64ToString(new64)

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
		panic(err)
	}

	logger.Debug("old pay id:", oldID)

	// update payinfo id
	oldID64, err := utils.StringToUint64(oldID)
	if err != nil {
		panic(err)
	}
	newID := utils.Uint64ToString(oldID64 + 1)
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
		Payer:  from,
		Credit: credit,
		TxHash: txHash,
	}
	// marshal to bytes
	data, err := json.Marshal(payInfo)
	if err != nil {
		panic(err)
	}
	// record payinfo data
	keys = append(keys, piKey)
	values = append(values, data)

	// multiput
	hc.LocalDB.MultiPut(keys, values)

	// todo: modify txinfo's state(credit saved)

	// response
	c.JSON(http.StatusOK, "pay ok")
}

// query pay infos
func (hc *HandlerCore) ListPayHandler(c *gin.Context) {
	addr := c.Query("addr")

	piList, err := hc.getPayInfoList(addr)
	if err != nil {
		panic(err)
	}

	// response
	c.JSON(http.StatusOK, piList)
}

// qeury credit for a role with address
func (hc *HandlerCore) QueryCreditHandler(c *gin.Context) {
	role := c.Query("role")
	address := c.Query("address")

	var credit string

	switch role {
	case "user":
		// get old credit from db, if key not found, init with 0
		credit, err := hc.getCredit(address)
		if err != nil {
			panic(err)
		}

		logger.Debug("credit:", credit)

		c.JSON(http.StatusOK, gin.H{
			credit: credit,
		})

	case "cp":
		// settle all orders of this cp, set order state

		// get cp's order list
		orderList, err := hc.getCpOrders(address)
		if err != nil {
			panic(err)
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
				panic(err)
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
					panic(err)
				}
				keys = append(keys, k)
				values = append(values, v)

				// set order's settled state to true
				k, v, err = hc.setOrderSettled([]byte(o.OrderKey), true)
				if err != nil {
					panic(err)
				}
				keys = append(keys, k)
				values = append(values, v)

				// multiput
				hc.LocalDB.MultiPut(keys, values)
			}
		}

		// get credit from db, if key not found, response 0
		credit, err = hc.getCredit(address)
		if err != nil {
			panic(err)
		}
		logger.Debug("credit:", credit)

		// response credit
		c.JSON(http.StatusOK, gin.H{
			"credit": credit,
		})
	default:
		c.JSON(http.StatusOK, gin.H{"res": "error role in request"})
	}
}

// transfer, write transfer records into db
// transfer info key: trans_user_id
// id key: trans_*
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
		panic(err)
	}
	id64, err := utils.StringToUint64(id)
	if err != nil {
		panic(err)
	}
	id64++
	newID := utils.Uint64ToString(id64)
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
		panic(err)
	}
	// record ti
	keys = append(keys, tiKey)
	values = append(values, data)

	// multiput
	hc.LocalDB.MultiPut(keys, values)

	c.JSON(http.StatusOK, gin.H{"res": "transfer ok"})
}

// list all transfer info about an user
func (hc *HandlerCore) ListTransferHandler(c *gin.Context) {

	// user address from param
	addr := c.Query("address")

	// transfer list for response
	transList, err := hc.getUserTransfers(addr)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, transList)
}

// refresh all transfers of an user, check transfers' confirmed state
func (hc *HandlerCore) RefreshTransferHandler(c *gin.Context) {
	userAddr := c.Query("address")

	transfers, err := hc.getUserTransfers(userAddr)
	if err != nil {
		panic(err)
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
			panic(err)
		}

		if confirmed {
			k, v, err := hc.setTransferConfirmed([]byte(transfer.TIKey), true)
			if err != nil {
				panic(err)
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
