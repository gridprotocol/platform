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

type PFServer struct {
	HttpServer *http.Server
}

// create new platform server with kv db
func NewServer(opt ServerOption) *PFServer {

	log.Println("Server Start")

	// init config
	err := config.InitConfig()
	if err != nil {
		log.Fatalf("failed to init the config: %v", err)
	}

	// make http server
	httSvr := &http.Server{
		Addr:    opt.Endpoint,
		Handler: routes.Routes{},
	}

	// make platform server
	pfServer := PFServer{
		HttpServer: httSvr,
	}
	return &pfServer
}

// register routes for http server
func (s *PFServer) RegisterRoutes() {
	// gin engine
	gin.SetMode(gin.ReleaseMode)

	// register routes
	router := routes.RegistRoutes()

	s.HttpServer.Handler = router
}
