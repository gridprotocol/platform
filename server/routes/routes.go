package routes

import (
	"github.com/gin-gonic/gin"
	_ "github.com/rockiecn/platform/docs"
	"github.com/rockiecn/platform/lib/kv"
	"github.com/rockiecn/platform/lib/log"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	logger         = log.Logger("routes")
	Chain_Endpoint string
)

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
	OrderKey string `json:"orderKey"`    // order key
	UserAddr string `json:"userAddress"` // user address
	CPAddr   string `json:"cpAddress"`   // provider address
	CPName   string `json:"cpName"`      // provider name
	EndPoint string `json:"endpoint"`    // provider endpoint
	NumCPU   string `json:"numCPU"`
	PriCPU   string `json:"priCPU"`
	NumGPU   string `json:"numGPU"`
	PriGPU   string `json:"priGPU"`
	NumDisk  string `json:"numDisk"`
	PriDisk  string `json:"priDisk"`
	NumMem   string `json:"numMem"`
	PriMem   string `json:"priMem"`
	Dur      string `json:"duration"`
	Expire   string `json:"expire"`
	Settled  bool   `json:"settled"`
	Cost     int64  `json:"cost"` // credit cost
}

// register all routes for server
func RegistRoutes(db *kv.Database, chain_ep string) Routes {

	// new default gin engine
	ginEng := gin.Default()

	// use cors middleware
	ginEng.Use(cors())

	routes := Routes{
		ginEng,
	}

	// store the chain endpoint for later use in hanlers
	Chain_Endpoint = chain_ep

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

	// for swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// for test
	r.GET("/", hc.RootHandler)

	// cp operation
	r.POST("/registcp", hc.RegistCPHandler)
	r.POST("/revisecp", hc.ReviseHandler)
	r.GET("/listcp", hc.ListCPHandler)
	r.GET("/getcp", hc.GetCPHandler)

	// query credit for an address
	r.GET("/querycredit", hc.QueryCreditHandler)

	// approve credit
	r.POST("/approve", hc.ApproveHandler)
	// check allowance after approve
	r.GET("allowance", hc.AllowanceHandler)

	// order operation
	r.POST("/createorder", hc.CreateOrderHandler)
	r.GET("/getorder", hc.GetOrderHandler)
	//r.GET("/listorder", hc.ListOrderHandler)

	r.POST("/userconfirm", hc.UserConfirmHandler)

}
