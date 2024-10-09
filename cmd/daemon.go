package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"runtime"
	"syscall"
	"time"

	"github.com/grid/contracts/eth"
	"github.com/mitchellh/go-homedir"
	"github.com/rockiecn/platform/common"
	"github.com/rockiecn/platform/lib/config"
	"github.com/rockiecn/platform/server"
	"github.com/urfave/cli/v2"
)

var DaemonCmd = &cli.Command{
	Name:  "daemon",
	Usage: "platform daemon",
	Subcommands: []*cli.Command{
		runCmd,
		stopCmd,
	},
}

// run daemon
var runCmd = &cli.Command{
	Name:  "run",
	Usage: "run server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "chain",
			Aliases: []string{"c"},
			Usage:   "chain to interactivate, local: use local test chain, sepo: use sepo test chain",
			Value:   "local",
		},
	},
	Action: func(ctx *cli.Context) error {
		chain := ctx.String("chain")

		// parse config file
		err := config.InitConfig()
		if err != nil {
			log.Fatalf("failed to init the config: %v", err)
		}
		ep := config.GetConfig().Http.Listen
		fmt.Println("endpoint:", ep)

		var chain_ep string

		// select contracts addresses for each chain
		switch chain {
		case "local":
			chain_ep = eth.Ganache
			common.Contracts = common.LocalContracts.Contracts

		case "sepo":
			chain_ep = eth.Sepolia
			common.Contracts = common.SepoContracts.Contracts
		}
		fmt.Printf("chain selected:%s, chain endpoint:%s\n", chain, chain_ep)
		fmt.Println("contract addresses:", common.Contracts)

		opts := server.ServerOption{
			Endpoint:       ep,
			Chain_Endpoint: chain_ep,
		}

		// create server
		svr := server.NewServer(opts)

		// open and close for server
		svr.LocalDB.Open()
		defer svr.LocalDB.Close()

		// register routes for server
		svr.RegisterRoutes()

		// start server
		go func() {
			if err := svr.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()

		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down server...")

		cctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := svr.HttpServer.Shutdown(cctx); err != nil {
			log.Fatal("Server forced to shutdown: ", err)
		}

		log.Println("Server exiting")
		return nil
	},
}

var stopCmd = &cli.Command{
	Name:  "stop",
	Usage: "stop server",
	Action: func(_ *cli.Context) error {
		pidpath, err := homedir.Expand("./")
		if err != nil {
			return nil
		}

		pd, _ := os.ReadFile(path.Join(pidpath, "pid"))

		err = kill(string(pd))
		if err != nil {
			return err
		}
		log.Println("gateway gracefully exit...")

		return nil
	},
}

func kill(pid string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("kill", "-15", pid).Run()
	case "windows":
		return exec.Command("taskkill", "/F", "/T", "/PID", pid).Run()
	default:
		return fmt.Errorf("unsupported platform %s", runtime.GOOS)
	}
}
