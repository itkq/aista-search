package controller

import (
	"aista-search/db"
	"aista-search/pagination"
	"aista-search/view"
	"github.com/gin-gonic/gin"
	"github.com/k0kubun/pp"
	"gopkg.in/guregu/null.v3"
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
	images, err = db.GetImages(
		null.NewInt(0, false),
		[]int{},
		null.StringFrom(q),
		true,
		false,
		null.NewInt(0, false),
	)
	if err != nil {
		pp.Println(err)
		c.String(500, "internal error")
		return
	}

	imagesVal := db.Images(*images)
	page, err := pagination.NewPagination(imagesVal.Interface(), p, db.ImagesPerPage)
	if err != nil {
		c.String(404, "404 page not found")
		return
	}

	v := view.New(c)
	v.Name = "search/index"
	v.Vars["q"] = q
	v.Vars["p"] = p
	v.Vars["page"] = page
	v.Render()
}
