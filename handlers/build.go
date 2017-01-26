package handlers

import (
	"github.com/autoabs/autoabs/build"
	"github.com/autoabs/autoabs/database"
	"github.com/autoabs/autoabs/utils"
	"github.com/gin-gonic/gin"
)

func buildGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	builds, err := build.GetAll(db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	build.Sort(builds)

	c.JSON(200, builds)
}

func buildArchive(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	buildId, ok := utils.ParseObjectId(c.Param("buildId"))
	if !ok {
		c.AbortWithStatus(400)
		return
	}

	bild, err := build.GetBuild(db, buildId)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	err = bild.Archive(db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, bild)
}
