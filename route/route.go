package route

import (
	"aista-search/controller"
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	router := gin.Default()

	router.GET("/", controller.IndexGET)

	router.GET("/episodes", controller.EpisodesGET)
	router.POST("/api/episode/create", controller.EpisodePOST)
	router.POST("/api/episode/update", controller.EpisodeUpdate)
	router.GET("/api/episode/latest", controller.LatestEpisodeGET)

	router.POST("/api/image/create", controller.ImagesPOST)

	return router
}
