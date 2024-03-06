package routes

import (
	"github.com/gin-gonic/gin"
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
	OrderID  string `json:"orderID"`     // order id for this user
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

func init() {

}

// register all routes for server
func RegistRoutes() Routes {

	router := gin.Default()

	router.Use(cors())

	r := Routes{
		router,
	}

	// register all routes
	r.registerAll()

	return r
}

// create local db, register all routes
func (r Routes) registerAll() {

	// new handler core
	hc := HandlerCore{
		LocalDB: nil,
	}
	// init db for hanlder core
	hc.InitDB()

	// for test
	r.GET("/", hc.RootHandler)

	// for functions
	r.POST("/registcp", hc.RegistCPHandler)
	r.GET("/listcp", hc.ListCPHandler)
	r.POST("/createorder", hc.CreateOrderHandler)
	// list orders for user
	r.GET("/listorder", hc.ListOrderHandler)
}
