package tender

import (
	"avito/internal/entity"

	"github.com/google/uuid"
)

type CreateTender struct {
	Name            string                   `json:"name" validate:"required,max=100"`
	Description     string                   `json:"description" validate:"required,max=500"`
	ServiceType     entity.TenderServiceType `json:"serviceType" validate:"required,oneof=Construction Delivery Manufacture"`
	OrganizationID  uuid.UUID                `json:"organizationId" validate:"required,max=100,uuid4"`
	CreatorUserName string                   `json:"creatorUsername" validate:"required" copier:"-"`
}

type PatchTender struct {
	Name        string                   `json:"name" validate:"max=100"`
	Description string                   `json:"description" validate:"max=500"`
	ServiceType entity.TenderServiceType `json:"serviceType" validate:"omitempty,oneof=Construction Delivery Manufacture"`
}
