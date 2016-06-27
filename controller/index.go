package controller

import (
	"github.com/gin-gonic/gin"
)

func IndexGET(c *gin.Context) {
	c.String(200, "Hello, World!")
}
