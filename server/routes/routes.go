package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rockiecn/platform/lib/config"
	"github.com/rockiecn/platform/lib/kv"
	"github.com/rockiecn/platform/lib/logs"
)

var logger = logs.Logger("local")

type Routes struct {
	*gin.Engine
}

type NodeInfo struct {
	Name     string `json:"name"`
	Entrance string `json:"entrance"`
	Resource string `json:"resource"`
	Price    string `json:"price"`
}

type OrderInfo struct {
	ID       string `json:"id"`
	Resource string `json:"resource"`
	Duration string `json:"duration"`
	Price    string `json:"price"`
}

type CPInfo struct {
	Addr  string `json:"address"`
	Name  string `json:"name"`
	Entr  string `json:"entrance"`
	Res   string `json:"resource"`
	Price string `json:"price"`
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
	// create kv db
	db, err := kv.NewDatabase(config.GetConfig().Local.DBPath)
	if err != nil {
		logger.Error("Fail to open up the database, err: ", err)
		panic(err)
	}

	// handler core
	hc := handlerCore{
		DB: db,
	}

	// for test
	r.GET("/", hc.RootHandler)

	// for functions
	r.GET("/list", hc.ListHandler)
	r.GET("/order", hc.OrderHandler)
	r.GET("/login", hc.LoginHandler)
}
