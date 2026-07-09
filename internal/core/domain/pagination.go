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
	if p.Limit != nil {
		if err := ValidateLimit(*p.Limit); err != nil {
			return err
		}
	}

	if p.Offset != nil {
		if err := ValidateOffset(*p.Offset); err != nil {
			return err
		}
	}

	return nil
}

func ValidateLimit(limit int) error {
	if limit < 0 {
		return ErrNegativeLimit
	}

	return nil
}

func ValidateOffset(offset int) error {
	if offset < 0 {
		return ErrNegativeOffset
	}

	return nil
}
