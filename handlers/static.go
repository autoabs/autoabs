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
		c.Writer.Header().Add("Expires", "0")
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

func staticGet(c *gin.Context) {
	staticPath(c, "/static"+c.Params.ByName("path"))
}

func staticTestingGet(c *gin.Context) {
	c.Writer.Header().Add("Cache-Control",
		"no-cache, no-store, must-revalidate")
	c.Writer.Header().Add("Pragma", "no-cache")
	c.Writer.Header().Add("Expires", "0")

	pth := c.Params.ByName("path")
	if pth == "" {
		if c.Request.URL.Path == "/config.js" {
			pth = "config.js"
		} else if c.Request.URL.Path == "/build.js" {
			pth = "build.js"
		} else {
			pth = "index.html"
		}
	}

	c.Writer.Header().Add("Content-Type", static.GetMimeType(pth))
	fileServer.ServeHTTP(c.Writer, c.Request)
}
