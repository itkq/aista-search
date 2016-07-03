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
	router.POST("/api/episode/update", controller.EpisodeUpdatePOST)
	router.GET("/api/episode/latest", controller.LatestEpisodeGET)

	router.GET("/api/images", controller.ImagesGET)
	router.POST("/api/image/create", controller.ImagesPOST)
	router.POST("/api/image/update", controller.ImagesUpdatePOST)

	return router
}
