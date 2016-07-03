package controller

import (
	"aista-search/db"
	"github.com/gin-gonic/gin"
	"github.com/k0kubun/pp"
	"strconv"
)

func ImagesPOST(c *gin.Context) {
	var images []db.Image
	if c.BindJSON(&images) != nil {
		c.JSON(400, gin.H{"status": "bad", "msg": "bad json format"})
		return
	}

	if err := db.CreateImages(images); err != nil {
		c.JSON(500, gin.H{"status": "bad", "msg": "db error"})
		return
	}

	c.JSON(200, gin.H{"status": "ok", "count": len(images)})
}

func ImagesUpdatePOST(c *gin.Context) {
	var images []db.Image
	if err := c.BindJSON(&images); err != nil {
		c.JSON(400, gin.H{"status": "bad", "msg": "bad json format"})
		return
	}

	if err := db.UpdateImages(images); err != nil {
		pp.Println(err)
		c.JSON(500, gin.H{"status": "bad", "msg": "db error"})
		return
	}

	c.JSON(200, gin.H{"status": "ok", "count": len(images)})
}

func ImagesGET(c *gin.Context) {
	episodeID, err := strconv.Atoi(c.Query("episode_id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "bad", "msg": "request error"})
		return
	}

	images, err := db.GetImagesByEpisodeID(episodeID)
	if err != nil {
		pp.Println(err)
		c.JSON(500, gin.H{"status": "bad", "msg": "db error"})
		return
	}

	c.JSON(200, gin.H{"status": "ok", "images": *images})
}
