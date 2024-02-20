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
func (hc *handlerCore) LoginCPHandler(c *gin.Context) {

	name := c.PostForm("name")

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

	// address as key, info as valude
	hc.CPDB.Put([]byte(name), []byte(data))

	c.JSON(http.StatusOK, "login OK")
}

// handler of create order
func (hc *handlerCore) CreateOrderHandler(c *gin.Context) {

	// user address
	addr := c.PostForm("address")

	// provider name
	name := c.PostForm("name")

	numCPU := c.PostForm("numCPU")
	priCPU := c.PostForm("priCPU")

	numGPU := c.PostForm("numGPU")
	priGPU := c.PostForm("priGPU")

	numStore := c.PostForm("numStore")
	priStore := c.PostForm("priStore")

	numMem := c.PostForm("numMem")
	priMem := c.PostForm("priMem")

	dur := c.PostForm("duration")

	// get order id for each user
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
	fmt.Println("old order id:", utils.BytesToInt(orderID))

	//
	info := OrderInfo{
		ID:       fmt.Sprintf("%d", utils.BytesToInt(orderID)),
		Addr:     addr,
		Name:     name,
		NumCPU:   numCPU,
		PriCPU:   priCPU,
		NumGPU:   numGPU,
		PriGPU:   priGPU,
		NumStore: numStore,
		PriStore: priStore,
		NumMem:   numMem,
		PriMem:   priMem,
		Dur:      dur,
	}

	// mashal into bytes
	data, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}

	// 'address' _ 'order id' as order key
	strKey := fmt.Sprintf("%s_%d", addr, utils.BytesToInt(orderID))
	fmt.Println("key:", strKey)

	// put order info into db
	hc.OrderDB.Put([]byte(strKey), []byte(data))

	// increase order id
	intID := utils.BytesToInt(orderID)
	intID += 1
	orderID = utils.IntToBytes(intID)
	fmt.Println("new order id:", utils.BytesToInt(orderID))
	// write new order id into db
	err = hc.OrderDB.Put([]byte(addr), orderID)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, "create order OK")
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
