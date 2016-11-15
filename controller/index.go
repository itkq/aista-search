package controller

import (
	"aista-search/view"
	"github.com/gin-gonic/gin"
)

func IndexGET(c *gin.Context) {
	v := view.New(c)
	v.Name = "index/index"
	v.Render()
}

func AboutGET(c *gin.Context) {
	v := view.New(c)
	v.Name = "index/about"
	v.Render()
}

func Ping(c *gin.Context) {
	c.String(200, "pong")
}
