package controller

import (
	"aista-search/db"
	"aista-search/session"
	"aista-search/view"
	"github.com/gin-gonic/gin"
	"github.com/k0kubun/pp"
	"gopkg.in/guregu/null.v3"
	"strconv"
)

func (a *api) ImagesGET(c *gin.Context) {
	episodeID, err := strconv.Atoi(c.Query("episode_id"))
	toUpload, _ := strconv.ParseBool(c.Query("to_upload"))
	cnt, _ := strconv.Atoi(c.Query("cnt"))

	images, err := db.GetImages(
		null.IntFrom(int64(episodeID)),
		null.NewString("", false),
		false,
		toUpload,
		null.IntFrom(int64(cnt)),
	)
	if err != nil {
		pp.Println(err)
		c.JSON(500, gin.H{"status": "bad", "msg": "db error"})
		return
	}

	c.JSON(200, gin.H{"status": "ok", "images": *images})
}

func (a *api) ImagesPOST(c *gin.Context) {
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

func (a *api) ImagesPUT(c *gin.Context) {
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

func ImagePOST(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	sentence := c.PostForm("sentence")
	sess := session.Instance(c.Request)

	var images []db.Image
	image, err := db.GetImageByID(id)
	if err != nil {
		sess.AddFlash(view.Flash{"更新エラーです", "error"})
		sess.Save(c.Request, c.Writer)
	} else {
		image.Sentence = null.StringFrom(sentence)
		images = append(images, *image)

		if err := db.UpdateImages(images); err != nil {
			sess.AddFlash(view.Flash{"更新エラーです", "error"})
			sess.Save(c.Request, c.Writer)
		} else {
			sess.AddFlash(view.Flash{"更新しました", "success"})
			sess.Save(c.Request, c.Writer)
		}
	}

	c.Redirect(302, "/images/"+strconv.Itoa(id))
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
