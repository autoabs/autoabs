package cmd

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/autoabs/autoabs/config"
	"github.com/autoabs/autoabs/handlers"
	"github.com/autoabs/autoabs/node"
	"github.com/autoabs/autoabs/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Starts web server node
func App() {
	var debug bool
	debugStr := os.Getenv("DEBUG")
	if debugStr == "" {
		debug = true
	} else {
		debug, _ = strconv.ParseBool(debugStr)
	}

	router := gin.New()

	if debug {
		router.Use(gin.Logger())
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	handlers.Register(router)

	nde := node.Node{
		Name: utils.RandName(),
		Type: "app",
	}
	nde.Keepalive()

	addr := fmt.Sprintf(
		"%s:%d",
		config.Config.ServerHost,
		config.Config.ServerPort,
	)

	server := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 4096,
	}

	logrus.WithFields(logrus.Fields{
		"address": config.Config.ServerHost,
		"port":    config.Config.ServerPort,
		"debug":   debug,
	}).Info("cmd.app: Starting app node")

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
