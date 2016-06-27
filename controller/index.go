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
