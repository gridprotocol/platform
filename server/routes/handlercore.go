package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rockiecn/platform/lib/kv"
)

type handlerCore struct {
	DB *kv.Database
}

// handler of welcom
func (hc *handlerCore) RootHandler(c *gin.Context) {
	c.String(http.StatusOK, "Welcome Server")
}

// handler for list cp nodes
func (hc *handlerCore) ListHandler(c *gin.Context) {

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

	allValue := hc.DB.GetAllValues()
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

// handler for get order
func (hc *handlerCore) OrderHandler(c *gin.Context) {
	or := OrderInfo{ID: "123", Resource: "res", Duration: "dur", Price: "100"}
	// response order
	c.JSON(http.StatusOK, or)

}

// handler of login
func (hc *handlerCore) LoginHandler(c *gin.Context) {
	addr := c.Query("address")
	name := c.Query("name")
	entrance := c.Query("entrance")
	resource := c.Query("resource")
	price := c.Query("price")
	info := CPInfo{
		Addr:  addr,
		Name:  name,
		Entr:  entrance,
		Res:   resource,
		Price: price,
	}

	// mashal into bytes
	data, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}

	// address as key, info as valude
	hc.DB.Put([]byte(addr), []byte(data))

	c.JSON(http.StatusOK, "login OK")
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
