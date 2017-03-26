package handlers

import (
	"github.com/autoabs/autoabs/database"
	"github.com/autoabs/autoabs/node"
	"github.com/gin-gonic/gin"
)

func nodeGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	builds, err := node.GetAll(db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, builds)
}
