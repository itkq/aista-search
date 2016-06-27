package controller

import (
	"aista-search/db"
	"github.com/gin-gonic/gin"
	"strconv"
)

func EpisodesGET(c *gin.Context) {
	episodes, err := db.GetEpisodes()

	if err != nil {
		c.JSON(400, nil)
		return
	}

	c.JSON(200, episodes)
}

func EpisodePOST(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	title := c.PostForm("title")

	if id == 0 || title == "" {
		c.JSON(400, gin.H{"status": "bad"})
		return
	}

	ep, err := db.CreateEpisode(id, title, 0)

	if err != nil {
		c.JSON(500, gin.H{"status": "bad"})
		return
	}

	c.JSON(200, gin.H{"status": "ok", "episode": ep})
}
