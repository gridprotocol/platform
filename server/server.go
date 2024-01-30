package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rockiecn/platform/lib/config"
	"github.com/rockiecn/platform/lib/logs"
	"github.com/rockiecn/platform/server/routes"
)

var logger = logs.Logger("server")

type ServerOption struct {
	Endpoint string
}

// create new platform server with kv db
func NewServer(opt ServerOption) *http.Server {

	log.Println("Server Start")

	// init config
	err := config.InitConfig()
	if err != nil {
		log.Fatalf("failed to init the config: %v", err)
	}

	// gin engine
	gin.SetMode(gin.ReleaseMode)
	// register routes
	router := routes.RegistRoutes()
	// http server
	svr := &http.Server{
		Addr:    opt.Endpoint,
		Handler: router,
	}

	return svr
}
