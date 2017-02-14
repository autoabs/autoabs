package handlers

import (
	"github.com/autoabs/autoabs/build"
	"github.com/autoabs/autoabs/database"
	"github.com/autoabs/autoabs/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

type buildData struct {
	Builds []*build.Build `json:"builds"`
	Index  int            `json:"index"`
	Count  int            `json:"count"`
}

func buildGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	indexStr := c.Query("index")
	index, _ := strconv.Atoi(indexStr)

	builds, index, count, err := build.GetAll(db, index)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	data := &buildData{
		Builds: builds,
		Index:  index,
		Count:  count,
	}

	c.JSON(200, data)
}

func buildLogGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	buildId, ok := utils.ParseObjectId(c.Param("buildId"))
	if !ok {
		c.AbortWithStatus(400)
		return
	}

	bild, err := build.Get(db, buildId)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, bild.Log)
}

func buildArchivePut(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	buildId, ok := utils.ParseObjectId(c.Param("buildId"))
	if !ok {
		c.AbortWithStatus(400)
		return
	}

	bild, err := build.Get(db, buildId)
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

func buildRebuildPut(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	buildId, ok := utils.ParseObjectId(c.Param("buildId"))
	if !ok {
		c.AbortWithStatus(400)
		return
	}

	bild, err := build.Get(db, buildId)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	err = bild.Rebuild(db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, bild)
}
