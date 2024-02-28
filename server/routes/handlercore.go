package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rockiecn/platform/lib/kv"
	"github.com/rockiecn/platform/lib/utils"
)

type handlerCore struct {
	CPDB    *kv.Database
	OrderDB *kv.Database
}

// handler of welcom
func (hc *handlerCore) RootHandler(c *gin.Context) {
	c.String(http.StatusOK, "Welcome Server")
}

// handler for list cp nodes
func (hc *handlerCore) ListCPHandler(c *gin.Context) {

	// // read db
	// data, err := hc.DB.Get([]byte("0x0090675FD3ef5031d7719A758163E73Fd58AF1EB"))
	// if err != nil {
	// 	panic(err)
	// }
	// logger.Info("data from db:", string(data))

	// // unmarshal
	// cpInfo := &CPInfo{}
	// err = json.Unmarshal(data, cpInfo)
	// if err != nil {
	// 	panic(err)
	// }
	// logger.Info("unmarshaled cp info:", cpInfo)

	// all cp info to response
	cps := make([]CPInfo, 0, 100)

	allValue := hc.CPDB.GetAllValues()
	for _, v := range allValue {
		cp := &CPInfo{}
		err := json.Unmarshal([]byte(v), cp)
		if err != nil {
			panic(err)
		}
		cps = append(cps, *cp)
	}

	// response node list
	c.JSON(http.StatusOK, cps)
}

// handler of cp login
func (hc *handlerCore) RegistCPHandler(c *gin.Context) {

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

	b, err := hc.CPDB.Has([]byte(address))
	if err != nil {
		panic(err)
	}
	if b {
		c.JSON(http.StatusOK, "cp already exist")
		return
	}

	// wallet address as key, info as valude
	err = hc.CPDB.Put([]byte(address), []byte(data))
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, "regist OK")
}

// handler of create order
func (hc *handlerCore) CreateOrderHandler(c *gin.Context) {

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

	// compute expire with duration and current time
	expire, err := utils.DurToTS(dur)
	if err != nil {
		panic(err)
	}

	// read cp name from db with cp address

	// check cp
	b, err := hc.CPDB.Has([]byte(cpAddr))
	if err != nil {
		panic(err)
	}
	if !b {
		c.JSON(http.StatusOK, "cp not found")
		return
	}
	// read cp info
	data, err := hc.CPDB.Get([]byte(cpAddr))
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
	orderID, err := hc.OrderDB.Get([]byte(userAddr))
	if err != nil {
		// if no order id, init with 0
		if err.Error() == "Key not found" {
			err = hc.OrderDB.Put([]byte(userAddr), utils.IntToBytes(0))
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	fmt.Println("old order id:", utils.BytesToInt(orderID))

	// construct new order info
	info := OrderInfo{
		OrderID:  fmt.Sprintf("%d", utils.BytesToInt(orderID)),
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
	}

	// mashal order info into bytes
	data, err = json.Marshal(info)
	if err != nil {
		panic(err)
	}

	// 'user address' _ 'order id' as order key
	strKey := fmt.Sprintf("%s_%d", userAddr, utils.BytesToInt(orderID))
	fmt.Println("key:", strKey)

	// put order info into db
	hc.OrderDB.Put([]byte(strKey), []byte(data))

	// increase order id
	intID := utils.BytesToInt(orderID)
	intID += 1
	orderID = utils.IntToBytes(intID)
	fmt.Println("new order id:", utils.BytesToInt(orderID))
	// update order id into db
	err = hc.OrderDB.Put([]byte(userAddr), orderID)
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

// handler for get order list
func (hc *handlerCore) ListOrderHandler(c *gin.Context) {

	// user address from param
	addr := c.Query("address")
	// get order id, equal to order number of this user
	orderID, err := hc.OrderDB.Get([]byte(addr))
	if err != nil {
		// if no order id, init with 0
		if err.Error() == "Key not found" {
			err = hc.OrderDB.Put([]byte(addr), utils.IntToBytes(0))
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	orderList := make([]OrderInfo, 0, 100)
	// number of order
	num := utils.BytesToInt(orderID)
	for i := 0; i < num; i++ {
		// make key
		key := fmt.Sprintf("%s_%d", addr, i)
		// get order
		data, err := hc.OrderDB.Get([]byte(key))
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

	// response order list
	c.JSON(http.StatusOK, orderList)

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
	nMem, err := utils.StringToUint64(o.NumMem)
	if err != nil {
		return 0, err
	}
	pMem, err := utils.StringToUint64(o.PriMem)
	if err != nil {
		return 0, err
	}
	nStor, err := utils.StringToUint64(o.NumStore)
	if err != nil {
		return 0, err
	}
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
