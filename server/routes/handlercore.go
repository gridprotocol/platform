package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/dgraph-io/badger/v2"
	"github.com/gin-gonic/gin"
	"github.com/rockiecn/platform/lib/kv"
	"github.com/rockiecn/platform/lib/utils"
)

type HandlerCore struct {
	LocalDB *kv.Database
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
			err := appendResult(&cps, it.Item())
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
func appendResult(cps *[]CPInfo, item *badger.Item) error {
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
	orderID, err := hc.LocalDB.Get([]byte(userAddr))
	if err != nil {
		// if no order id, init with 0
		if err.Error() == "Key not found" {
			err = hc.LocalDB.Put([]byte(userAddr), utils.IntToBytes(0))
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	logger.Debugf("old order id:", utils.BytesToInt(orderID))

	// 'user address' _ 'order id' as order key
	orderKey := fmt.Sprintf("%s_%d", userAddr, utils.BytesToInt(orderID))
	logger.Debugf("key:", orderKey)

	// construct new order info
	info := OrderInfo{
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

	// mashal order info into bytes
	data, err = json.Marshal(info)
	if err != nil {
		panic(err)
	}

	// put order info into db
	hc.LocalDB.Put([]byte(orderKey), []byte(data))

	// increase order id
	intID := utils.BytesToInt(orderID)
	intID += 1
	orderID = utils.IntToBytes(intID)
	logger.Debugf("new order id:", utils.BytesToInt(orderID))
	// update order id
	err = hc.LocalDB.Put([]byte(userAddr), orderID)
	if err != nil {
		panic(err)
	}

	// append an order key for cp
	err = hc.appendOrder(cpAddr, orderKey)
	if err != nil {
		panic(err)
	}

	// calc value of order
	v, err := CalcValue(&info)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"to":    "0x1234",
		"value": v,
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
		orderList, err = hc.listUserOrder(addr)
		if err != nil {
			panic(err)
		}
	case "cp":
		orderList, err = hc.listCpOrder(addr)
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
func (hc *HandlerCore) listUserOrder(userAddr string) ([]OrderInfo, error) {
	orderList := make([]OrderInfo, 0, 100)

	// get order id, equal to order number of this user
	orderID, err := hc.LocalDB.Get([]byte(userAddr))
	if err != nil {
		// if no order id, init with 0
		if err.Error() == "Key not found" {
			return orderList, nil
		} else {
			return nil, err
		}
	}

	// number of order
	num := utils.BytesToInt(orderID)
	for i := 0; i < num; i++ {
		// make key
		key := fmt.Sprintf("%s_%d", userAddr, i)
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
func (hc *HandlerCore) listCpOrder(cpAddr string) ([]OrderInfo, error) {
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
func (hc *HandlerCore) appendOrder(cpAddr string, orderKey string) error {
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

	// put new order list for cp
	err = hc.LocalDB.Put([]byte(cpordersKey), data)
	if err != nil {
		panic(err)
	}

	return nil
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

// calc value of an order
func CalcValue(o *OrderInfo) (uint64, error) {
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

	// get value
	value := (nCPU*pCPU + nGPU*pGPU + nMem*pMem + nStor*pStor) * dur

	return value, nil
}

// check number
func isNumber(s string) bool {
	pattern := `^[0-9]+(\.[0-9]+)?$`
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}
