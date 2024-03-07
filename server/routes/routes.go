package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rockiecn/platform/lib/kv"
	"github.com/rockiecn/platform/lib/logs"
)

var logger = logs.Logger("routes")

type Routes struct {
	*gin.Engine
}

type NodeInfo struct {
	Name     string `json:"name"`
	Entrance string `json:"entrance"`
	Resource string `json:"resource"`
	Price    string `json:"price"`
}

// type OrderInfo struct {
// 	ID       string `json:"id"`
// 	Resource string `json:"resource"`
// 	Duration string `json:"duration"`
// 	Price    string `json:"price"`
// }

type CPInfo struct {
	Name     string `json:"name"` // provider name
	Address  string `json:"address"`
	EndPoint string `json:"endpoint"`
	NumCPU   string `json:"numCPU"`
	PriCPU   string `json:"priCPU"`
	NumGPU   string `json:"numGPU"`
	PriGPU   string `json:"priGPU"`
	NumStore string `json:"numStore"`
	PriStore string `json:"priStore"`
	NumMem   string `json:"numMem"`
	PriMem   string `json:"priMem"`
}

type OrderInfo struct {
	OrderKey string `json:"orderKey"`    // order id for this user
	UserAddr string `json:"userAddress"` // user address
	CPAddr   string `json:"cpAddress"`   // provider address
	CPName   string `json:"cpName"`      // provider name
	EndPoint string `json:"endpoint"`    // provider endpoint
	NumCPU   string `json:"numCPU"`
	PriCPU   string `json:"priCPU"`
	NumGPU   string `json:"numGPU"`
	PriGPU   string `json:"priGPU"`
	NumStore string `json:"numStore"`
	PriStore string `json:"priStore"`
	NumMem   string `json:"numMem"`
	PriMem   string `json:"priMem"`
	Dur      string `json:"duration"`
	Expire   string `json:"expire"`
	Settled  bool   `json:"settled"`
}

// register all routes for server
func RegistRoutes(db *kv.Database) Routes {

	ginEng := gin.Default()

	ginEng.Use(cors())

	routes := Routes{
		ginEng,
	}

	// register all routes
	routes.registerAll(db)

	return routes
}

// create local db, register all routes
func (r Routes) registerAll(db *kv.Database) {

	// new handler core with db
	hc := HandlerCore{
		LocalDB: db,
	}

	// for test
	r.GET("/", hc.RootHandler)

	// cp operation
	r.POST("/registcp", hc.RegistCPHandler)
	r.GET("/listcp", hc.ListCPHandler)

	// order operation
	r.POST("/createorder", hc.CreateOrderHandler)
	// list orders for user
	r.GET("/listorder", hc.ListOrderHandler)

	// recharge credit with eth in tx
	r.POST("/credit", hc.Credit)
	// query credit for an address
	r.GET("/querycredit", hc.QueryCredit)
}
