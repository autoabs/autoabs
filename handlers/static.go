package handlers

import (
	"fmt"
	"github.com/autoabs/autoabs/constants"
	"github.com/gin-gonic/gin"
	"strings"
)

func staticPath(c *gin.Context, path string) {
	path = constants.StaticRoot + path

	fmt.Println(path)

	file, ok := store.Files[path]
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
