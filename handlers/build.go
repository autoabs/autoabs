package handlers

import (
	"github.com/autoabs/autoabs/build"
	"github.com/autoabs/autoabs/database"
	"github.com/gin-gonic/gin"
)

func buildGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	builds, err := build.GetAll(db)
	if err != nil {
		return
	}

	build.Sort(builds)

	c.JSON(200, builds)
}
