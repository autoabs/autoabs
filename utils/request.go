package utils

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func GetRemoteAddr(c *gin.Context) (addr string) {
	addr = c.Request.Header.Get("CF-Connecting-IP")
	if addr != "" {
		return
	}

	addr = c.Request.Header.Get("X-Forwarded-For")
	if addr != "" {
		return
	}

	addr = c.Request.Header.Get("X-Real-Ip")
	if addr != "" {
		return
	}

	addr = strings.Split(c.Request.RemoteAddr, ":")[0]
	return
}
