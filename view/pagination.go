package view

type Pagination struct {
	First     int
	PrevExist bool
	Prev      int
	Current   int
	Next      int
	NextExist bool
	Last      int
	List      []interface{}
}

func NewPagination(list []interface{}, current int, perPage int) (*Pagination, error) {
	var sliceList []interface{}

	offset := (current - 1) * perPage
	listSize := len(list)

	if offset == listSize {
		return nil, newPagingError("List is empty")
	}

	if listSize-offset <= perPage {
		sliceList = list[offset:listSize]
	} else {
		sliceList = list[offset : offset+perPage]
	}

	first := 1
	last := (listSize + perPage - 1) / perPage

	var prev, next int
	var prevExist, nextExist bool

	if current > 2 {
		prev = current - 1
	} else {
		prev = 0
	}

	if current+2 <= last {
		next = current + 1
	} else {
		next = 0
	}

	if current >= 4 {
		prevExist = true
	} else {
		prevExist = false
	}

	if current+3 <= last {
		nextExist = true
	} else {
		nextExist = false
	}

	if first == current {
		current = first
		first = 0
	} else if last == current {
		current = last
		last = 0
	} else if prev == current || next == current {
		current = 0
	}

	if last == 1 {
		last = 0
	}

	return &Pagination{
		First:     first,
		PrevExist: prevExist,
		Prev:      prev,
		Current:   current,
		Next:      next,
		NextExist: nextExist,
		Last:      last,
		List:      sliceList,
	}, nil
}

type PagingError struct {
	msg string
}

func (err *PagingError) Error() string {
	return err.msg
}

func newPagingError(s string) *PagingError {
	err := new(PagingError)
	err.msg = s
	return err
}
