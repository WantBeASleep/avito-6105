package parsers

import (
	"avito/api/validation"
	"avito/internal/entity"
	"net/http"
	"strconv"
)

const (
	DefaultLimit  = 5
	DefaultOffset = 0
)

func ParsePagination(r *http.Request) (*entity.Pagination, error) {
	params := entity.Pagination{
		Limit:  DefaultLimit,
		Offset: DefaultOffset,
	}

	if limit := r.URL.Query().Get("limit"); limit != "" {
		parsedLimit, err := strconv.Atoi(limit)
		if err != nil || parsedLimit < 0 {
			return nil, validation.NewValidateError("limit must be simple positive num")
		}
		params.Limit = parsedLimit
	}

	if offset := r.URL.Query().Get("offset"); offset != "" {
		parsedOffset, err := strconv.Atoi(offset)
		if err != nil || parsedOffset < 0 {
			return nil, validation.NewValidateError("offset must be simple positive num")
		}
		params.Offset = parsedOffset
	}

	return &params, nil
}
