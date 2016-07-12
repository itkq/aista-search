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

	router.GET("/image/:id", controller.ImageGET)
	router.POST("/image/update", controller.ImageUpdatePOST)
	router.GET("/api/images", controller.ImagesGET)
	router.GET("/api/images/upload", controller.ImagesToUploadGET)
	router.POST("/api/image/create", controller.ImagesPOST)
	router.POST("/api/image/update", controller.ImagesUpdatePOST)

	router.GET("/search", controller.SearchGET)
	api := router.Group("/api")
	{
		api.GET("/episodes/", controller.API.EpisodesGET)
		api.GET("/episodes/:id", controller.API.EpisodeGET)
		api.POST("/episodes/", controller.API.EpisodePOST)
		api.PUT("/episodes/:id", controller.API.EpisodePUT)
	}

	return router
}
