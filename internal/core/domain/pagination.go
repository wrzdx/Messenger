package domain

type Pagination struct {
	Limit  *int
	Offset *int
}

func NewPagination(limit, offset *int) Pagination {
	return Pagination{
		Limit:  limit,
		Offset: offset,
	}
}

func (p Pagination) Validate() error {
	if p.Limit != nil && *p.Limit < 0 {
		return ErrNegativeLimit
	}

	if p.Offset != nil && *p.Offset < 0 {
		return ErrNegativeOffset
	}

	return nil
}
