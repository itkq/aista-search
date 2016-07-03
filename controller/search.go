package controller

import (
	"aista-search/db"
	"aista-search/view"
	"github.com/gin-gonic/gin"
	"github.com/k0kubun/pp"
)

func SearchGET(c *gin.Context) {
	var images *[]db.Image
	var err error

	q := c.Query("q")
	if q == "" {
		images, err = db.GetImages()
		if err != nil {
			pp.Println(err)
			c.String(500, "internal error")
			return
		}
	} else {
		images, err = db.GetImagesBySentence(q)
		if err != nil {
			pp.Println(err)
			c.String(500, "internal error")
			return
		}
	}

	v := view.New(c)
	v.Name = "search/index"
	v.Vars["q"] = q
	v.Vars["images"] = *images
	v.Render()
}
