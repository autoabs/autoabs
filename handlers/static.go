package handlers

import (
	"github.com/autoabs/autoabs/constants"
	"github.com/autoabs/autoabs/static"
	"github.com/gin-gonic/gin"
	"strings"
)

func staticPath(c *gin.Context, pth string) {
	pth = constants.StaticRoot + pth

	file, ok := store.Files[pth]
	if !ok {
		c.AbortWithStatus(404)
		return
	}

	if constants.StaticCache {
		c.Writer.Header().Add("Cache-Control", "public, max-age=86400")
		c.Writer.Header().Add("ETag", file.Hash)
	} else {
		c.Writer.Header().Add("Cache-Control",
			"no-cache, no-store, must-revalidate")
		c.Writer.Header().Add("Pragma", "no-cache")
	}

	if strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
		c.Writer.Header().Add("Content-Encoding", "gzip")
		c.Data(200, file.Type, file.GzipData)
	} else {
		c.Data(200, file.Type, file.Data)
	}
}

func staticIndexGet(c *gin.Context) {
	staticPath(c, "/index.html")
}

func staticAppGet(c *gin.Context) {
	staticPath(c, "/app"+c.Params.ByName("path"))
}

func staticStylesGet(c *gin.Context) {
	staticPath(c, "/styles"+c.Params.ByName("path"))
}

func staticVendorGet(c *gin.Context) {
	staticPath(c, "/vendor"+c.Params.ByName("path"))
}

func staticLiveGet(c *gin.Context) {
	c.Writer.Header().Add("Cache-Control",
		"no-cache, no-store, must-revalidate")
	c.Writer.Header().Add("Pragma", "no-cache")
	c.Writer.Header().Add("Expires", "0")

	pth := c.Params.ByName("path")
	if pth == "" {
		pth = "index.html"
	}

	c.Writer.Header().Add("Content-Type", static.GetMimeType(pth))
	fileServer.ServeHTTP(c.Writer, c.Request)
}
