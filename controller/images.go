package controller

import (
	"aista-search/db"
	"aista-search/view"
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

func ImageGET(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	image, err := db.GetImageByID(id)
	if err != nil {
		c.String(404, "not found")
		return
	}

	v := view.New(c)
	v.Name = "images/detail"
	v.Vars["image"] = image
	v.Render()
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
