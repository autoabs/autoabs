package cmd

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/autoabs/autoabs/config"
	"github.com/autoabs/autoabs/constants"
	"github.com/autoabs/autoabs/handlers"
	"github.com/autoabs/autoabs/node"
	"github.com/autoabs/autoabs/scheduler"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func WebNode() {
	if constants.Production {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	if !constants.Production {
		router.Use(gin.Logger())
	}

	handlers.Register(router)

	nde := node.Node{
		Id:   config.Config.WebNodeId,
		Type: "web",
	}
	nde.Init()

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
		"address":    config.Config.ServerHost,
		"port":       config.Config.ServerPort,
		"production": constants.Production,
	}).Info("cmd.app: Starting app node")

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func StorageNode() {
	nde := node.Node{
		Id:   config.Config.StorageNodeId,
		Type: "storage",
	}
	nde.Init()

	sch := scheduler.Storage{}

	sch.Start()
}

func BuilderNode() {
	nde := node.Node{
		Id:   config.Config.BuilderNodeId,
		Type: "builder",
	}
	nde.Init()

	sch := scheduler.Build{}

	sch.Start()
}
