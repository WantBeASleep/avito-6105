package parsers

import (
	"avito/api/validation"
	"avito/internal/entity"
	"fmt"
	"net/http"
)

var AvailableServiceTypes = map[string]struct{}{
	"Construction": {},
	"Delivery":     {},
	"Manufacture":  {},
}

func ParseServiceTypes(r *http.Request) ([]entity.TenderServiceType, error) {
	var serviceTypes []entity.TenderServiceType
	parsedServiceTypes := r.URL.Query()["service_type"]
	for _, s := range parsedServiceTypes {
		if _, ok := AvailableServiceTypes[s]; !ok {
			return nil, validation.NewValidateError(fmt.Sprintf("invalid service type value: %s", s))
		}
		serviceTypes = append(serviceTypes, entity.TenderServiceType(s))
	}

	return serviceTypes, nil
}
