package pagination

import (
	"aista-search/db"
	"github.com/k0kubun/pp"
	"testing"
)

type pagination struct {
	First     int
	PrevExist bool
	Prev      int
	Current   int
	Next      int
	NextExist bool
	Last      int
}

type testingWrapper testing.T

func TestPagination(t *testing.T) {
	var images db.Images
	var rowPage *Pagination
	var page, expected pagination
	tw := testingWrapper(*t)

	// Normal list
	images = newImages(100)
	rowPage, _ = NewPagination(images.Interface(), 2, db.ImagesPerPage)
	tw.assertSize(rowPage.List, db.ImagesPerPage)

	page = convert(*rowPage)
	expected = pagination{
		First:     1,
		PrevExist: false,
		Prev:      0,
		Current:   2,
		Next:      3,
		NextExist: false,
		Last:      4,
	}
	if page != expected {
		pp.Println(page)
		t.Error("pagination error")
	}

	// Empty list
	images = newImages(0)
	rowPage, err := NewPagination(images.Interface(), 1, db.ImagesPerPage)
	if err == nil {
		t.Error("pagination error")
	}

	// Empty list for paging
	images = newImages(db.ImagesPerPage)
	rowPage, err = NewPagination(images.Interface(), 2, db.ImagesPerPage)
	if err.Error() != "List is empty" {
		t.Error("pagination error")
	}

	// Single page
	images = newImages(29)
	rowPage, _ = NewPagination(images.Interface(), 1, db.ImagesPerPage)
	tw.assertSize(rowPage.List, db.ImagesPerPage)

	page = convert(*rowPage)
	expected = pagination{
		First:     0,
		PrevExist: false,
		Prev:      0,
		Current:   1,
		Next:      0,
		NextExist: false,
		Last:      0,
	}
	if page != expected {
		pp.Println(page)
		t.Error("pagination error")
	}

	// Pattern 1
	images = newImages(140)
	rowPage, _ = NewPagination(images.Interface(), 4, db.ImagesPerPage)
	tw.assertSize(rowPage.List, db.ImagesPerPage)

	page = convert(*rowPage)
	expected = pagination{
		First:     1,
		PrevExist: true,
		Prev:      3,
		Current:   4,
		Next:      0,
		NextExist: false,
		Last:      5,
	}
	if page != expected {
		pp.Println(page)
		t.Error("pagination error")
	}

	// Pattern 2
	images = newImages(140)
	rowPage, _ = NewPagination(images.Interface(), 3, db.ImagesPerPage)
	tw.assertSize(rowPage.List, db.ImagesPerPage)

	page = convert(*rowPage)
	expected = pagination{
		First:     1,
		PrevExist: false,
		Prev:      2,
		Current:   3,
		Next:      4,
		NextExist: false,
		Last:      5,
	}
	if page != expected {
		pp.Println(page)
		t.Error("pagination error")
	}

	// Pattern 3
	images = newImages(160)
	rowPage, _ = NewPagination(images.Interface(), 2, db.ImagesPerPage)
	tw.assertSize(rowPage.List, db.ImagesPerPage)

	page = convert(*rowPage)
	expected = pagination{
		First:     1,
		PrevExist: false,
		Prev:      0,
		Current:   2,
		Next:      3,
		NextExist: true,
		Last:      6,
	}
	if page != expected {
		pp.Println(page)
		t.Error("pagination error")
	}

	// Pattern 4
	images = newImages(40)
	rowPage, _ = NewPagination(images.Interface(), 2, db.ImagesPerPage)
	tw.assertSize(rowPage.List, db.ImagesPerPage)

	page = convert(*rowPage)
	expected = pagination{
		First:     1,
		PrevExist: false,
		Prev:      0,
		Current:   2,
		Next:      0,
		NextExist: false,
		Last:      0,
	}
	if page != expected {
		pp.Println(page)
		t.Error("pagination error")
	}
}

func (tw *testingWrapper) assertSize(list []interface{}, size int) {
	if len(list) != size {
		tw.Error("size error")
	}
}

func convert(p Pagination) pagination {
	return pagination{
		First:     p.First,
		PrevExist: p.PrevExist,
		Prev:      p.Prev,
		Current:   p.Current,
		Next:      p.Next,
		NextExist: p.NextExist,
		Last:      p.Last,
	}
}

func newImages(size int) db.Images {
	return db.Images(make([]db.Image, size))
}
