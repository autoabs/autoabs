package handlers

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/autoabs/autoabs/database"
	"github.com/autoabs/autoabs/utils"
	"github.com/dropbox/godropbox/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Limit size of request body
func Limiter(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1000000)
}

// Get database from session
func Database(c *gin.Context) {
	db := database.GetDatabase()
	c.Set("db", db)
	c.Next()
	db.Close()
}

// Recover panics
func Recovery(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithFields(logrus.Fields{
				"client": utils.GetRemoteAddr(c),
				"error":  errors.New(fmt.Sprintf("%s", r)),
			}).Error("handlers: Handler panic")
			c.Writer.WriteHeader(http.StatusInternalServerError)
		}
	}()

	c.Next()
}

// Register all endpoint handlers
func Register(engine *gin.Engine) {
	engine.Use(Limiter)
	engine.Use(Recovery)

	dbGroup := engine.Group("")
	dbGroup.Use(Database)

	engine.GET("/check", checkGet)

	dbGroup.GET("/builds", buildsGet)
}
