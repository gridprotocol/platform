package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rockiecn/platform/lib/config"
	"github.com/rockiecn/platform/lib/kv"
	"github.com/rockiecn/platform/lib/log"
	"github.com/rockiecn/platform/server/routes"
)

var logger = log.Logger("server")

type ServerOption struct {
	Endpoint       string
	Chain_Endpoint string
}

type PFServer struct {
	LocalDB        *kv.Database
	HttpServer     *http.Server
	Chain_Endpoint string
}

// create new platform server with kv db
func NewServer(opt ServerOption) *PFServer {

	logger.Info("Server Start")

	// init config
	err := config.InitConfig()
	if err != nil {
		logger.Fatalf("failed to init the config: %v", err)
	}

	// make http server
	httSvr := &http.Server{
		Addr:    opt.Endpoint,
		Handler: routes.Routes{},
	}

	// make platform server
	pfServer := PFServer{
		HttpServer:     httSvr,
		LocalDB:        nil,
		Chain_Endpoint: opt.Chain_Endpoint,
	}

	// init db
	pfServer.InitDB()

	return &pfServer
}

// register routes for http server
func (s *PFServer) RegisterRoutes() {
	// gin engine
	gin.SetMode(gin.ReleaseMode)

	// register routes
	routes := routes.RegistRoutes(s.LocalDB, s.Chain_Endpoint)

	s.HttpServer.Handler = routes
}

// init db for platform
func (s *PFServer) InitDB() {
	// get db path from config
	db_path := config.GetConfig().Local.DB_Path
	// create cp db
	s.LocalDB = kv.NewDB(db_path)
}
