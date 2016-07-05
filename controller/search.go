package controller

import (
	"aista-search/db"
	"aista-search/pagination"
	"aista-search/view"
	"github.com/gin-gonic/gin"
	"github.com/k0kubun/pp"
	"strconv"
)

func SearchGET(c *gin.Context) {
	var images *[]db.Image
	var err error

	p, _ := strconv.Atoi(c.Query("p"))
	if p == 0 {
		p = 1
	}

	q := c.Query("q")
	if q == "" {
		images, err = db.GetImages(p)
		if err != nil {
			pp.Println(err)
			c.String(500, "internal error")
			return
		}
	} else {
		images, err = db.GetImagesBySentence(q, p)
		if err != nil {
			pp.Println(err)
			c.String(500, "internal error")
			return
		}
	}

	imagesVal := db.Images(*images)
	page, err := pagination.NewPagination(imagesVal.Interface(), p, db.ImagesPerPage)
	if err != nil {
		c.String(500, "paging error")
		return
	}

	v := view.New(c)
	v.Name = "search/index"
	v.Vars["q"] = q
	v.Vars["p"] = p
	v.Vars["page"] = page
	v.Render()
}
