package controller

import (
	"aista-search/db"
	"aista-search/view"
	"github.com/gin-gonic/gin"
	"strconv"
)

func EpisodesGET(c *gin.Context) {
	episodes, _ := db.GetEpisodes()

	v := view.New(c)
	v.Name = "episodes/index"
	v.Vars["episodes"] = episodes
	v.Render()
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

	c.JSON(200, gin.H{"status": "ok", "id": ep.Id})
}
