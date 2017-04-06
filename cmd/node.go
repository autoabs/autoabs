package cmd

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/autoabs/autoabs/config"
	"github.com/autoabs/autoabs/handlers"
	"github.com/autoabs/autoabs/node"
	"github.com/autoabs/autoabs/scheduler"
	"github.com/autoabs/autoabs/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"time"
)

func WebNode() {
	var debug bool
	debugStr := os.Getenv("DEBUG")
	if debugStr == "" {
		debug = true
	} else {
		debug, _ = strconv.ParseBool(debugStr)
	}

	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	if debug {
		router.Use(gin.Logger())
	}

	handlers.Register(router)

	nde := node.Node{
		Id:   utils.RandName(),
		Type: "web",
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

func StorageNode() {
	nde := node.Node{
		Id:   utils.RandName(),
		Type: "storage",
	}
	nde.Keepalive()

	sch := scheduler.Storage{}

	sch.Start()
}

func BuilderNode() {
	nde := node.Node{
		Id:   utils.RandName(),
		Type: "builder",
	}
	nde.Keepalive()

	sch := scheduler.Build{}

	sch.Start()
}
