package route

import (
	"aista-search/controller"
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	router := gin.Default()

	router.GET("/", controller.IndexGET)

	return router
}
