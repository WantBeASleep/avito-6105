package parsers

import (
	"avito/api/validation"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var (
	ParserEmptyString = func(s string) (string, error) { return s, nil }

	ParserUUID = func(s string) (uuid.UUID, error) {
		parse, err := uuid.Parse(s)
		if err != nil {
			return uuid.Nil, err
		}
		return parse, nil
	}

	ParserInt = func(s string) (int, error) {
		parse, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return parse, nil
	}
)

func parse[T any](value string, parseName string, requiredFlag bool, parser func(string) (T, error)) (T, error) {
	var emptyParam T

	if value == "" {
		if requiredFlag {
			return emptyParam, validation.NewValidateError(fmt.Sprintf("%s required", parseName))
		}
		return emptyParam, nil
	}

	parseParam, err := parser(value)
	if err != nil {
		return emptyParam, validation.NewValidateError(fmt.Sprintf("invalid %s format", parseName))
	}

	return parseParam, nil
}

func ParseVar[T any](r *http.Request, parseName string, requiredFlag bool, parser func(string) (T, error)) (T, error) {
	paramStr := mux.Vars(r)[parseName]
	return parse(paramStr, parseName, requiredFlag, parser)
}

func ParseQuery[T any](r *http.Request, parseName string, requiredFlag bool, parser func(string) (T, error)) (T, error) {
	paramStr := r.URL.Query().Get(parseName)
	return parse(paramStr, parseName, requiredFlag, parser)
}
