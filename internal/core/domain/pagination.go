package domain

import "errors"

var (
	ErrNegativeLimit  = errors.New("limit must be non-negative")
	ErrNegativeOffset = errors.New("offset must be non-negative")
)

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
	fields := make(map[string]string)
	if p.Limit != nil {
		if err := ValidateLimit(*p.Limit); err != nil {
			fields["limit"] = err.Error()
		}
	}

	if p.Offset != nil {
		if err := ValidateOffset(*p.Offset); err != nil {
			fields["offset"] = err.Error()
		}
	}

	if len(fields) > 0 {
		return ValidationErr(PaginationEntity, fields)
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
