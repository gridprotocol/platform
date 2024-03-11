package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
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
	PIKey       string `json:"PayInfoKey"`
	Payer       string `json:"Payer"`
	Credit      uint64 `json:"Credit"`
	TxHash      string `json:"TxHash"`
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

	// mashal into bytes
	data, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}

	// cp key: cp_info_*
	key := fmt.Sprintf("cp_info_%s", address)

	// check if cp exists
	b, err := hc.LocalDB.Has([]byte(key))
	if err != nil {
		panic(err)
	}
	if b {
		c.JSON(http.StatusOK, "cp already exist")
		return
	}

	// wallet address as key, info as valude
	err = hc.LocalDB.Put([]byte(key), []byte(data))
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, "regist OK")
}

// handler for list cp nodes
func (hc *HandlerCore) ListCPHandler(c *gin.Context) {

	// all cp info to response
	cps := make([]CPInfo, 0, 100)

	prefix := []byte("cp_info_") // 设置通配符前缀
	err := hc.LocalDB.DB.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
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

// append db item into cps
func (hc *HandlerCore) appendResult(cps *[]CPInfo, item *badger.Item) error {
	// append each item
	err := item.Value(func(val []byte) error {
		logger.Debugf("Key:%s Value:%s", string(item.Key()), string(val))
		cp := CPInfo{}
		err := json.Unmarshal(val, &cp)
		if err != nil {
			return err
		}
		// append
		*cps = append(*cps, cp)
		return nil
	})
	if err != nil {
		return fmt.Errorf("error processing result: %w", err)
	}

	return nil
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

	cpkey := fmt.Sprintf("cp_info_%s", cpAddr)
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
	data, err = hc.LocalDB.Get([]byte(userAddr))
	if err != nil {
		// if no order id, init with 0
		if err.Error() == "Key not found" {
			orderID = "0"
		} else {
			panic(err)
		}
	} else {
		logger.Debugf("data:%s", data)
		orderID = string(data)
	}

	logger.Debugf("old order id:%s", orderID)

	// 'user address' _ 'order id' as order key
	orderKey := fmt.Sprintf("%s_%s", userAddr, orderID)
	logger.Debugf("order key:%s", orderKey)

	// construct new order info
	order := OrderInfo{
		OrderKey: orderKey,
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

	// check credit for user
	var credit string
	creKey := fmt.Sprintf("cre_%s", userAddr)
	data, err = hc.LocalDB.Get([]byte(creKey))
	if err != nil {
		if err.Error() == "Key not found" {
			c.JSON(http.StatusOK, "no credit for this user")
			return
		} else {
			panic(err)
		}
	} else {
		credit = string(data)
	}
	logger.Debug("credit left:", credit)
	credit64, err := utils.StringToUint64(credit)
	if err != nil {
		panic(err)
	}
	// check credit
	if credit64 < cost64 {
		c.JSON(http.StatusOK, "credit is not enough to pay this order,create order failed")
		return
	}

	// for atomic operations on db
	keys := [][]byte{}
	values := [][]byte{}
	pos := 0

	// update user's credit
	credit64 = credit64 - cost64
	credit = utils.Uint64ToString(credit64)
	// update user's credit in db
	//hc.LocalDB.Put([]byte(creKey), []byte(credit))
	keys = append(keys, []byte(creKey))
	values = append(values, []byte(credit))
	pos++
	logger.Debug("new credit after createorder:", credit)

	// db operation

	// mashal order into bytes
	data, err = json.Marshal(order)
	if err != nil {
		panic(err)
	}
	// put order info into db
	hc.LocalDB.Put([]byte(orderKey), []byte(data))

	// increase order id by 1
	orderID64, err := utils.StringToUint64(orderID)
	if err != nil {
		panic(err)
	}
	orderID64 += 1
	orderID = utils.Uint64ToString(orderID64)
	logger.Debugf("new order id:%s", orderID)
	// update order id
	// err = hc.LocalDB.Put([]byte(userAddr), []byte(orderID))
	// if err != nil {
	// 	panic(err)
	// }
	keys = append(keys, []byte(userAddr))
	values = append(values, []byte(orderID))
	pos++

	// append the order key for cp into db
	k, v, err := hc.appendOrder(cpAddr, orderKey)
	if err != nil {
		panic(err)
	}
	keys = append(keys, k)
	values = append(values, v)
	pos++

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

	// order list for response
	orderList := make([]OrderInfo, 0, 100)

	var err error

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

// get user's order list from db
func (hc *HandlerCore) getUserOrders(userAddr string) ([]OrderInfo, error) {
	orderList := make([]OrderInfo, 0, 100)

	// get order id, equal to order number of this user
	data, err := hc.LocalDB.Get([]byte(userAddr))
	if err != nil {
		// if no order id, init with 0
		if err.Error() == "Key not found" {
			return orderList, nil
		} else {
			return nil, err
		}
	}
	orderID := string(data)
	logger.Debug("user's order id:", orderID)

	// number of order
	num, err := utils.StringToUint64(orderID)
	if err != nil {
		panic(err)
	}
	for i := uint64(0); i < num; i++ {
		// make key
		key := fmt.Sprintf("%s_%d", userAddr, i)
		logger.Debug("order key:", key)
		// get order
		data, err := hc.LocalDB.Get([]byte(key))
		if err != nil {
			return nil, err
		}
		order := &OrderInfo{}
		err = json.Unmarshal(data, order)
		if err != nil {
			return nil, err
		}
		orderList = append(orderList, *order)
	}

	return orderList, nil
}

// get order list for cp
func (hc *HandlerCore) getCpOrders(cpAddr string) ([]OrderInfo, error) {
	// 'cp' _ 'address' as cp key
	cpordersKey := fmt.Sprintf("cp_orders_%s", cpAddr)

	// init an empty order list
	orderList := make([]OrderInfo, 0, 100)

	// read db for cp order keys data
	data, err := hc.LocalDB.Get([]byte(cpordersKey))
	if err != nil {
		// if no order id, return empty order list
		if err.Error() == "Key not found" {
			return orderList, nil
		} else {
			return nil, err
		}
	}

	var orderKeys []string
	// unmarshal data into order keys if data is not empty
	if len(data) != 0 {
		err = json.Unmarshal([]byte(data), &orderKeys)
		if err != nil {
			return nil, err
		}
	} else { // if no key data, return empty list
		return orderList, nil
	}

	// get order list with order keys
	for i := 0; i < len(orderKeys); i++ {
		// each item is an order key
		key := orderKeys[i]
		// get order
		data, err := hc.LocalDB.Get([]byte(key))
		if err != nil {
			panic(err)
		}
		order := &OrderInfo{}
		err = json.Unmarshal(data, order)
		if err != nil {
			panic(err)
		}
		orderList = append(orderList, *order)
	}

	return orderList, nil
}

// append an order key for a cp
func (hc *HandlerCore) appendOrder(cpAddr string, orderKey string) (k []byte, v []byte, err error) {
	// 'cp' _ 'address' as cp key
	cpordersKey := fmt.Sprintf("cp_orders_%s", cpAddr)

	var orderKeys []string = make([]string, 0)

	// read order keys from db
	data, err := hc.LocalDB.Get([]byte(cpordersKey))
	if err != nil {
		// if no order keys, init an empty data
		if err.Error() == "Key not found" {
			data = []byte{}
		} else {
			panic(err)
		}
	}

	// if data not empty, unmarshal it
	if len(data) != 0 {
		err = json.Unmarshal(data, &orderKeys)
		if err != nil {
			panic(err)
		}
	}

	// append into keys
	orderKeys = append(orderKeys, orderKey)

	data, err = json.Marshal(orderKeys)
	if err != nil {
		panic(err)
	}

	// return new order list for cp
	// err = hc.LocalDB.Put([]byte(cpordersKey), data)
	// if err != nil {
	// 	panic(err)
	// }

	return []byte(cpordersKey), data, nil
}

// user record credit with txHash
// value uint: eth
// crecit = eth * 1000000
func (hc *HandlerCore) Pay(c *gin.Context) {
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

	// credit key: cre_*
	creKey := fmt.Sprintf("cre_%s", from)
	logger.Debug("credit key:", creKey)

	// update credit for this account
	var old string
	// get old credit from db, if key not found, init with 0
	data, err := hc.LocalDB.Get([]byte(creKey))
	if err != nil {
		if err.Error() == "Key not found" {
			old = "0"
		} else {
			panic(err)
		}
	} else {
		old = string(data)
	}
	logger.Debug("old credit:", old)
	logger.Debug("credit:", credit)

	// accumulate credit
	old64, err := utils.StringToUint64(old)
	if err != nil {
		panic(err)
	}
	new64 := old64 + credit
	new := utils.Uint64ToString(new64)

	logger.Debug("new credit:", new)

	// update credit for this account
	err = hc.LocalDB.Put([]byte(creKey), []byte(new))
	if err != nil {
		panic(err)
	}

	// get payinfo id for this account: payid_*
	var oldID string
	idKey := fmt.Sprintf("payid_%s", from)
	data, err = hc.LocalDB.Get([]byte(idKey))
	if err != nil {
		if err.Error() == "Key not found" {
			oldID = "0"
		} else {
			panic(err)
		}
	} else {
		oldID = string(data)
	}
	logger.Debug("old pay id:", oldID)

	oldID64, err := utils.StringToUint64(oldID)
	if err != nil {
		panic(err)
	}

	// update payinfo id
	newID := utils.Uint64ToString(oldID64 + 1)
	logger.Debug("new payinfo id:", newID)
	// update payinfo id for this account
	err = hc.LocalDB.Put([]byte(idKey), []byte(newID))
	if err != nil {
		panic(err)
	}

	// make payinfo's key: account_pi_id
	piKey := fmt.Sprintf("%s_pi_%s", from, oldID)
	logger.Debug("payinfo key:", piKey)
	// record pay info into db
	// payInfo: from, value, credit, txHash, txConfirmed, creditSaved
	payInfo := PayInfo{
		PIKey:       piKey,
		Payer:       from,
		Credit:      credit,
		TxHash:      txHash,
		TxConfirmed: false,
		CreditSaved: false,
	}
	// marshal to bytes
	data, err = json.Marshal(payInfo)
	if err != nil {
		panic(err)
	}
	// record payinfo data
	err = hc.LocalDB.Put([]byte(piKey), data)
	if err != nil {
		panic(err)
	}

	// response
	c.JSON(http.StatusOK, "pay ok")
}

// query pay infos
func (hc *HandlerCore) QueryPay(c *gin.Context) {
	addr := c.Query("addr")

	piList, err := hc.getPayInfoList(addr)
	if err != nil {
		panic(err)
	}

	// response
	c.JSON(http.StatusOK, piList)
}

// get user's order list from db
func (hc *HandlerCore) getPayInfoList(addr string) ([]PayInfo, error) {
	piList := make([]PayInfo, 0, 100)

	idKey := fmt.Sprintf("payid_%s", addr)
	// get pay id
	data, err := hc.LocalDB.Get([]byte(idKey))
	if err != nil {
		// if no pay id, no payinfo
		if err.Error() == "Key not found" {
			return piList, nil
		} else {
			return nil, err
		}
	}
	payID := string(data)
	logger.Debug("account's pay id:", payID)

	// number of order
	num, err := utils.StringToUint64(payID)
	if err != nil {
		panic(err)
	}
	for i := uint64(0); i < num; i++ {
		// make payInfo key
		piKey := fmt.Sprintf("%s_pi_%d", addr, i)
		logger.Debug("order key:", piKey)
		// get order
		data, err := hc.LocalDB.Get([]byte(piKey))
		if err != nil {
			return nil, err
		}
		pi := &PayInfo{}
		err = json.Unmarshal(data, pi)
		if err != nil {
			return nil, err
		}
		piList = append(piList, *pi)
	}

	// payinfo list
	return piList, nil
}

// qeury credit for a role with address
func (hc *HandlerCore) QueryCredit(c *gin.Context) {
	role := c.Query("role")
	address := c.Query("address")

	// key: cre_address
	creKey := fmt.Sprintf("cre_%s", address)

	var credit string

	switch role {
	case "user":
		// get old credit from db, if key not found, init with 0
		data, err := hc.LocalDB.Get([]byte(creKey))
		if err != nil {
			if err.Error() == "Key not found" {
				c.JSON(http.StatusOK, "account not found")
				return
			} else {
				panic(err)
			}
		} else {
			credit = string(data)
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
			pos := 0

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
				pos++
				// set order's settled state to true
				k, v, err = hc.setOrderSettled([]byte(o.OrderKey), true)
				if err != nil {
					panic(err)
				}
				keys = append(keys, k)
				values = append(values, v)
				pos++

				hc.LocalDB.MultiPut(keys, values)
			}
		}

		// get credit from db, if key not found, response 0
		data, err := hc.LocalDB.Get([]byte(creKey))
		if err != nil {
			if err.Error() == "Key not found" {
				c.JSON(http.StatusOK, gin.H{
					"credit": "0",
				})
				return
			} else {
				panic(err)
			}
		} else {
			credit = string(data)
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

// add an account with some credit, return k,v for db write
func (hc *HandlerCore) addCredit(addr string, credit uint64) (k []byte, v []byte, err error) {

	// credit key: cre_*
	creKey := fmt.Sprintf("cre_%s", addr)

	// old credit
	var old string

	// get old credit from db, if key not found, init with 0
	data, err := hc.LocalDB.Get([]byte(creKey))
	if err != nil {
		if err.Error() == "Key not found" {
			old = "0"
		} else {
			return nil, nil, err
		}
	} else {
		old = string(data)
	}
	logger.Debug("old credit:", old)
	logger.Debug("credit:", credit)

	// accumulate credit
	old64, err := utils.StringToUint64(old)
	if err != nil {
		return nil, nil, err
	}
	new64 := old64 + credit
	new := utils.Uint64ToString(new64)

	logger.Debug("new credit:", new)

	// update credit for this account
	// err = hc.LocalDB.Put([]byte(creKey), []byte(new))
	// if err != nil {
	// 	return nil, nil, err
	// }

	// return k v for multiput
	return []byte(creKey), []byte(new), nil
}

// set an order's settled state with key
func (hc *HandlerCore) setOrderSettled(key []byte, settled bool) (k []byte, v []byte, err error) {
	// get order info with key
	data, err := hc.LocalDB.Get([]byte(key))
	if err != nil {
		return nil, nil, err
	}

	order := OrderInfo{}
	err = json.Unmarshal(data, &order)
	if err != nil {
		return nil, nil, err
	}

	// set order's settled state
	order.Settled = settled

	// marshal new order
	data, err = json.Marshal(order)
	if err != nil {
		return nil, nil, err
	}

	// update order info
	// err = hc.LocalDB.Put(key, data)
	// if err != nil {
	// 	return nil, nil, err
	// }

	// return k,v for db
	return key, data, nil
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

// calc credit cost of an order
func CalcCost(o *OrderInfo) (uint64, error) {
	nCPU, err := utils.StringToUint64(o.NumCPU)
	if err != nil {
		return 0, err
	}
	pCPU, err := utils.StringToUint64(o.PriCPU)
	if err != nil {
		return 0, err
	}

	nGPU, err := utils.StringToUint64(o.NumGPU)
	if err != nil {
		return 0, err
	}
	pGPU, err := utils.StringToUint64(o.PriGPU)
	if err != nil {
		return 0, err
	}

	// T
	nMem, err := utils.StringToUint64(o.NumMem)
	if err != nil {
		return 0, err
	}
	// T to byte
	nMem = nMem * 1024 * 1024 * 1024 * 1024
	pMem, err := utils.StringToUint64(o.PriMem)
	if err != nil {
		return 0, err
	}

	// G
	nStor, err := utils.StringToUint64(o.NumStore)
	if err != nil {
		return 0, err
	}
	// G to byte
	nStor = nStor * 1024 * 1024 * 1024
	pStor, err := utils.StringToUint64(o.PriStore)
	if err != nil {
		return 0, err
	}

	dur, err := utils.StringToUint64(o.Dur)
	if err != nil {
		return 0, err
	}

	// get wei value
	value := (nCPU*pCPU + nGPU*pGPU + nMem*pMem + nStor*pStor) * dur
	logger.Debug("wei of order:", value)

	cost := value / 1000 / 1000 / 1000 / 1000
	logger.Debug("credit cost of order:", cost)

	// return credit cost
	return cost, nil
}

// check number
func isNumber(s string) bool {
	pattern := `^[0-9]+(\.[0-9]+)?$`
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}
