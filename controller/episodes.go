package controller

import (
	"aista-search/db"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func (a *api) EpisodeGET(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		c.JSON(400, gin.H{"status": "bad", "msg": "invalid query"})
		return
	}

	episode, err := db.GetEpisode(id)
	if err != nil {
		c.JSON(500, gin.H{"status": "bad", "msg": "db error"})
		log.Println(err)
		return
	}

	if episode.ID == 0 {
		c.JSON(400, gin.H{"status": "bad", "msg": "no episode"})
		return
	}

	c.JSON(200, gin.H{"status": "ok", "episode": episode, "msg": ""})
}

func (a *api) EpisodesGET(c *gin.Context) {
	var cnt int
	var err error

	if c.Query("cnt") == "" {
		cnt = 10
	} else {
		cnt, err = strconv.Atoi(c.Query("cnt"))
	}

	if err != nil || cnt == 0 {
		c.JSON(400, gin.H{"status": "bad", "msg": "invalid query"})
		return
	}

	episodes, err := db.GetEpisodes(cnt)
	if err != nil {
		c.JSON(500, gin.H{"status": "bad", "msg": "db error"})
		log.Println(err)
		return
	}

	c.JSON(200, gin.H{"status": "ok", "episodes": episodes})
}

func (a *api) EpisodePOST(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	title := c.PostForm("title")

	if id == 0 || title == "" {
		c.JSON(400, gin.H{"status": "bad", "msg": "invalid request"})
		return
	}

	ep, err := db.CreateEpisode(id, title, db.EpCreated)

	if err != nil {
		c.JSON(500, gin.H{"status": "bad", "msg": "db error"})
		log.Println(err)
		return
	}

	c.JSON(200, gin.H{"status": "ok", "id": ep.ID, "msg": ""})
}

func (a *api) EpisodePUT(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	title := c.PostForm("title")
	status, _ := strconv.ParseUint(c.PostForm("status"), 10, 16)

	if id == 0 {
		c.JSON(400, gin.H{"status": "bad", "msg": "invalid request"})
		return
	}

	if _, err := db.GetEpisode(id); err != nil {
		c.JSON(500, gin.H{"status": "bad", "msg": "no episode"})
		log.Println(err)
		return
	}

	err = db.UpdateEpisode(id, title, uint(status))
	if err != nil {
		c.JSON(500, gin.H{"status": "bad", "msg": "db error"})
		log.Println(err)
		return
	}

	c.JSON(200, gin.H{"status": "ok", "msg": ""})
}
