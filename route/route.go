package route

import (
	"aista-search/config"
	"aista-search/controller"
	"github.com/gin-gonic/gin"
	"os"
)

func New() *gin.Engine {
	if os.Getenv("GO_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	router.Static("/img", config.GetEnv("IMG_ROOT", "./img"))

	router.GET("/", controller.IndexGET)
	router.GET("/ping", controller.Ping)

	router.GET("/images/:id", controller.ImageGET)
	router.POST("/images/:id", controller.ImagePOST)

	router.GET("/search", controller.SearchGET)

	api := router.Group("/api")
	{
		api.GET("/episodes/", controller.API.EpisodesGET)
		api.GET("/episodes/:id", controller.API.EpisodeGET)
		api.POST("/episodes/", controller.API.EpisodePOST)
		api.PUT("/episodes/:id", controller.API.EpisodePUT)

		api.GET("/images/", controller.API.ImagesGET)
		api.POST("/images/", controller.API.ImagesPOST)
		api.PUT("/images/", controller.API.ImagesPUT)
	}

	return router
}
