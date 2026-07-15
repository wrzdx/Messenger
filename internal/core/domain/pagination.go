package domain

import "errors"

var (
	ErrInvalidPagination = errors.New("invalid pagination options")
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
	if p.Limit != nil && *p.Limit < 0 {
		fields["limit"] = "limit must be non-negative"
	}

	if p.Offset != nil && *p.Offset < 0 {
		fields["offset"] = "offset must be non-negative"
	}

	if len(fields) > 0 {
		return DetailedError{
			Err:     ErrInvalidPagination,
			Details: fields,
		}
	}

	return nil
}
