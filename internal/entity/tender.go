package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TenderStatusType string

const (
	Created   TenderStatusType = "Created"
	Published TenderStatusType = "Published"
	Closed    TenderStatusType = "Closed"
)

var TenderStatusTypeList = []TenderStatusType{Created, Published, Closed}

type TenderServiceType string

const (
	Construction TenderServiceType = "Construction"
	Delivery     TenderServiceType = "Delivery"
	Manufacture  TenderServiceType = "Manufacture"
)

var TenderServiceTypeList = []TenderServiceType{Construction, Delivery, Manufacture}

type Tender struct {
	Id             uuid.UUID         `json:"id"`
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	Status         TenderStatusType  `json:"status"`
	ServiceType    TenderServiceType `json:"serviceType"`
	OrganizationID uuid.UUID         `json:"organizationId"`
	Version        int               `json:"version"`
	CreatedAt      time.Time         `json:"createdAt"`
}

func (t Tender) MarshalJSON() ([]byte, error) {
	type Alias Tender
	return json.Marshal(
		struct {
			*Alias
			CreatedAt string `json:"createdAt"`
		}{
			Alias:     (*Alias)(&t),
			CreatedAt: t.CreatedAt.Format(time.RFC3339),
		},
	)
}
