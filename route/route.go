package route

import (
	"aista-search/config"
	"aista-search/controller"
	"aista-search/route/middleware"
	"github.com/gin-gonic/gin"
	"os"
)

func New() *gin.Engine {
	if os.Getenv("GO_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	router.Static("/img", config.GetEnv("IMG_ROOT", "./img"))
	router.Static("/js", config.GetEnv("JS_ROOT", "./assets/js"))
	router.Static("/css", config.GetEnv("JS_ROOT", "./assets/css"))

	router.GET("/", controller.SearchGET)
	router.GET("/about", controller.AboutGET)
	router.GET("/ping", controller.Ping)

	router.GET("/images/:id", controller.ImageGET)
	router.POST("/images/:id", controller.ImagePOST)
	router.GET("/admin/images/", controller.ImagesGET)
	router.POST("/admin/images/", controller.ImagesDELETE)

	api := router.Group("/api")
	api.Use(middleware.Auth())
	{
		api.GET("/episodes/", controller.API.EpisodesGET)
		api.GET("/episodes/:id", controller.API.EpisodeGET)
		api.POST("/episodes/", controller.API.EpisodePOST)
		api.PUT("/episodes/:id", controller.API.EpisodePUT)

		api.GET("/images/", controller.API.ImagesGET)
		api.POST("/images/", controller.API.ImagesPOST)
		api.PUT("/images/", controller.API.ImagesPUT)
		api.DELETE("/images/:id", controller.API.ImageDelete)
	}

	return router
}
