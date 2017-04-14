package handlers

import (
	"github.com/autoabs/autoabs/database"
	"github.com/autoabs/autoabs/node"
	"github.com/gin-gonic/gin"
)

func nodeGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	nodes, err := node.GetAll(db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, nodes)
}

func nodePut(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	nodeId := c.Param("nodeId")
	if nodeId == "" {
		c.AbortWithStatus(400)
		return
	}

	nde, err := node.Get(db, nodeId)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	err = c.Bind(nde.Settings)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	err = nde.CommitSetttings(db)
	if err != nil {
		return
	}

	c.JSON(200, nde)
}
